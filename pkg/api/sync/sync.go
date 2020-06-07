// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package sync manages the exchange of subscriptions, slices and grains between primary and
// secondary servers. It contains different endpoints for primary and secondary servers
package sync

/*
  The Primary adds subscriptions (slices assigned to companies) and the Secondary asks for
  those assigned to them. This means that the Secondary (who currently initiates the sync)
  begins a sync session by asking for their subscriptions. It then adds any new ones to their
  local database.

  There is also an "active" subscription flag that can be changed (on either side). If disabled
  on the Primary, it will update the Secondary and log the activity. If enabled on the Primary,
  it will not change the Secondary. The active flag on the Secondary controls if it tries to sync
  that subscription, but changes are not propagated to the Primary. So, all of this means that
  the Primary controls what can be synced, but the Secondary can turn the sync off.

  The sync process will also observe the "active" company flag (on both sides) and the "allow_sync"
  slice is being updated flag (on the Primary).
*/

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/shared/client"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
)

// todo: change sync process to a websocket connection instead of separate http calls

type subsArray []sandpiper.Subscription

// Start sends a sync request to a primary sandpiper server from our secondary server
func (s *Sync) Start(c echo.Context, primaryID uuid.UUID) (err error) {
	var p *sandpiper.Company

	// log activity even if early exit
	defer func(begin time.Time) {
		msg := fmt.Sprintf("Syncing \"%s\" (%s)", p.Name, p.SyncAddr)
		if e := s.sdb.LogActivity(s.db, primaryID, uuid.Nil, msg, time.Since(begin), err); e != nil {
			err = fmt.Errorf("%w; LogActivity Error: %v", err, e)
		}
	}(time.Now())

	// must be a secondary server to start the sync
	if err := s.rbac.EnforceServerRole(sandpiper.SecondaryServer); err != nil {
		return err
	}
	// must be a local admin to start the sync
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return err
	}
	// get company information for the primary server we're syncing
	p, err = s.sdb.Primary(s.db, primaryID)
	if err != nil {
		return err
	}
	// connect to the primary server using their api-key (saving token)
	s.api, err = s.connect(p.SyncAddr, p.SyncAPIKey)
	if err != nil {
		return err
	}
	// get our subscriptions (with slices) from the primary server
	primSubs, err := s.api.AllSubs()
	if err != nil {
		return err
	}
	// get local subscriptions (with slices) as a receiver for this primary company
	localSubs, err := s.sdb.Subscriptions(s.db, primaryID)
	if err != nil {
		return nil
	}
	// sync all active subscriptions
	if err := s.syncSubscriptions(localSubs, primSubs, primaryID); err != nil {
		return err
	}
	return nil
}

func (s *Sync) connect(addr, key string) (*client.Client, error) {
	server, err := url.ParseRequestURI(addr)
	if err != nil {
		return nil, err
	}
	if key == "" {
		return nil, errors.New("api-key is missing")
	}
	api, err := client.SyncLogin(server, key, false) // nowhere to get a debug flag
	if err != nil {
		return nil, err
	}
	return api, nil
}

// syncSubscriptions makes sure our local subscriptions match primary ones. If we don't
// have a subscription, add it locally. If disabled on the Primary, disable it on the
// Secondary and log the activity. If enabled on the Primary but not on the secondary,
// do not make any changes. Perform a grain sync on all unlocked active slices.
func (s *Sync) syncSubscriptions(locals, prims subsArray, primaryID uuid.UUID) (err error) {
	// save our local subscriptions in a dictionary
	subs := make(sandpiper.SubsMap)
	subs.Load(locals)

	// run through primary (remote) subs and add to ours (local) if missing
	// (the slice for a new subscription is also added, but not the slice metadata, which
	// is added during the sync process)
	for _, remote := range prims {
		local, found := subs[remote.SubID]
		if !found {
			// add this subscription (and its slice) locally
			local = remote.SemiDeepCopy()
			local.CompanyID = primaryID  // change to our frame of reference for the add
			local.Slice.ContentHash = "" // force a re-sync
			if err := s.sdb.AddSlice(s.db, local.Slice); err != nil {
				return err
			}
			if err := s.sdb.AddSubscription(s.db, local); err != nil {
				return err
			}
		} else {
			// see if we should deactivate our active subscription (and so not process it)
			if !remote.Active && local.Active {
				if err := s.sdb.DeactivateSubscription(s.db, local.SubID); err != nil {
					return err
				}
				local.Active = false
			}
		}
		if local.Active {
			// sync the grains for a slice
			if err := s.syncSlice(primaryID, local.SubID, local.Slice, remote.Slice); err != nil {
				return err
			}
		}
	}
	return nil
}

// syncSlice does the actual work of looking for changes and doing the update.
// surround the sync with a begin/finalize update of the slice row to approximate a transaction
// and log results (and errors) to the activity table
func (s *Sync) syncSlice(primaryID, subID uuid.UUID, localSlice, remoteSlice *sandpiper.Slice) (err error) {
	// log activity at slice level *only* if an error occurs
	defer func(begin time.Time) {
		duration := time.Since(begin)
		msg := "Slice \"" + localSlice.Name + "\""
		if err != nil {
			if e := s.sdb.LogActivity(s.db, primaryID, subID, msg, duration, err); e != nil {
				err = fmt.Errorf("%w; LogActivity Error: %v", err, e)
			}
		}
		// log every sync attempt to primary (ignoring error)
		_ = s.api.LogActivity(s.rbac.OurServer().ID, subID, msg, duration, err)
	}(time.Now())

	if !remoteSlice.AllowSync {
		return errors.New("slice locked (being updated) on server")
	}

	if slicesMatch(remoteSlice, localSlice) {
		// nothing to do
		return nil
	}

	// get remote grain list of ids
	remoteIDs, err := s.api.GrainIDs(remoteSlice.ID)
	if err != nil {
		return err
	}

	// get local grain list of just the ids
	localIDs, err := s.sdb.Grains(s.db, localSlice.ID, true)
	if err != nil {
		return err
	}

	// determine local changes required to make local grains match remote grains
	adds, deletes := compareSlices(remoteIDs, localIDs)

	// start syncing grains (with a quasi-transaction on the slice row)
	if err := s.sdb.BeginSyncUpdate(s.db, remoteSlice.ID); err != nil {
		return err
	}
	defer func(sliceID uuid.UUID) {
		if e := s.sdb.FinalizeSyncUpdate(s.db, sliceID, err); e != nil {
			err = fmt.Errorf("FinalizeSyncUpdate Error: %v (%w)", e, err)
		}
	}(remoteSlice.ID)

	// remove obsolete grains (if any)
	if err := s.sdb.DeleteGrains(s.db, deletes); err != nil {
		return err
	}

	// add new grains
	for _, grainID := range adds {
		grain, err := s.api.Grain(grainID)
		if err != nil {
			return err
		}
		if err := s.sdb.AddGrain(s.db, grain); err != nil {
			return err
		}
	}

	// replace local slice metadata with remote's
	meta, err := s.api.SliceMetaData(remoteSlice.ID)
	if err != nil {
		return err
	}
	if err := s.sdb.ReplaceSliceMetadata(s.db, remoteSlice.ID, meta); err != nil {
		return err
	}
	// Update ContentHash, ContentCount & ContentDate and verify our own hash against remote's
	err = s.sdb.RefreshSlice(s.db, remoteSlice)

	return err
}

// slicesMatch checks if slice has chanced and so needs to be updated
func slicesMatch(remoteSlice, localSlice *sandpiper.Slice) bool {
	// we can safely use the previous hash saved for comparison because we performed a deep
	// check of our own content when completing the previous sync
	return remoteSlice.ContentHash == localSlice.ContentHash
}

// compareSlices returns adds and deletes necessary to make the secondary match primary
func compareSlices(primary, secondary []sandpiper.Grain) (adds []uuid.UUID, dels []uuid.UUID) {
	// load primary grain ids into a set (marked as not matched)
	p := make(map[uuid.UUID]bool)
	for _, grain := range primary {
		p[grain.ID] = false
	}
	// create a list of secondary grain ids not in primary (marking matches)
	for _, grain := range secondary {
		if _, ok := p[grain.ID]; ok {
			p[grain.ID] = true
		} else {
			dels = append(dels, grain.ID)
		}
	}
	// all ids in primary not matched (still false) should be added
	for k, matched := range p {
		if !matched {
			adds = append(adds, k)
		}
	}
	return adds, dels
}

// Subscriptions returns all subscriptions with slices and metadata (not paginated)
// for the current user's company
func (s *Sync) Subscriptions(c echo.Context) ([]sandpiper.Subscription, error) {
	if err := s.rbac.EnforceServerRole(sandpiper.PrimaryServer); err != nil {
		return nil, err
	}
	if err := s.rbac.EnforceRole(c, sandpiper.SyncRole); err != nil {
		return nil, err
	}
	companyID := s.rbac.CurrentUser(c).CompanyID
	return s.sdb.Subscriptions(s.db, companyID)
}

// Grains returns all grains for a slice without pagination (with option to limit fields returned)
// Too bad we need to check company access to this slice again, but this is a public endpoint
// with no state beyond the user token. At least it uses a unique key for the check.
// Websockets should allow a more efficient approach.
func (s *Sync) Grains(c echo.Context, sliceID uuid.UUID, briefFlag bool) ([]sandpiper.Grain, error) {
	if err := s.rbac.EnforceServerRole(sandpiper.PrimaryServer); err != nil {
		return nil, err
	}
	if err := s.rbac.EnforceRole(c, sandpiper.SyncRole); err != nil {
		return nil, err
	}
	companyID := s.rbac.CurrentUser(c).CompanyID
	if err := s.sdb.SliceAccess(s.db, companyID, sliceID); err != nil {
		return nil, err
	}
	return s.sdb.Grains(s.db, sliceID, briefFlag)
}

// Process (NOT CURRENTLY USED) responds to a sync start request and "upgrades" http to a websocket
func (s *Sync) Process(c echo.Context) error {
	if err := s.rbac.EnforceServerRole(sandpiper.PrimaryServer); err != nil {
		return err
	}
	if err := s.rbac.EnforceRole(c, sandpiper.SyncRole); err != nil {
		return err
	}

	/*
		var (
			upgrader = websocket.Upgrader{}
		)

		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer ws.Close()

		// message loop
		for {
			// Write
			err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
			if err != nil {
				c.Logger().Error(err)
			}

			// Read
			_, msg, err := ws.ReadMessage()
			if err != nil {
				c.Logger().Error(err)
			}
			fmt.Printf("%s\n", msg)
		}
	*/
	return nil
}
