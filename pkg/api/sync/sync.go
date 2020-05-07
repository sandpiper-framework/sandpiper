// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package sync contains services for syncs. Syncs must belong to a slice
// and do not have an update method (use add/delete).
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
  slice updating flag (on the Primary).
*/

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"sandpiper/pkg/api/sync/credentials"
	"sandpiper/pkg/shared/client"
	"sandpiper/pkg/shared/model"
)

type subsArray []sandpiper.Subscription

// Start sends a sync request to a primary sandpiper server from our secondary server
func (s *Sync) Start(c echo.Context, primaryID uuid.UUID) (err error) {
	var p *sandpiper.Company

	// must be a secondary server to start the sync
	if err := s.rbac.EnforceServerRole("secondary"); err != nil {
		return err
	}
	// must be a local admin to start the sync
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return err
	}

	// log activity even if an error occurs
	defer func(begin time.Time) {
		msg := fmt.Sprintf("Syncing \"%s\" (%s)", p.Name, p.SyncAddr)
		_ = s.sdb.LogActivity(s.db, uuid.Nil, msg, time.Since(begin), err)
	}(time.Now())

	// get company information for the primary server we're syncing
	p, err = s.sdb.Primary(s.db, primaryID)
	if err != nil {
		return err
	}
	// connect to the primary server (saving token)
	api, err := s.connect(p.SyncAddr, p.SyncAPIKey, s.sec.APIKeySecret())
	if err != nil {
		return err
	}

	// todo: change sync process to a websocket connection instead of separate http calls

	// get our subscriptions (with slices) from the primary server
	primSubs, err := api.AllSubs()
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

func (s *Sync) connect(addr, key, secret string) (*client.Client, error) {
	// get login credentials from company's sync_api_key
	creds, err := credentials.New(key, secret)
	if err != nil {
		return nil, err
	}
	server, err := url.ParseRequestURI(addr)
	if err != nil {
		return nil, err
	}
	api, err := client.Login(server, creds.User, creds.Password, false)
	if err != nil {
		return nil, err
	}
	return api, nil
}

// syncSubscriptions makes sure our local subscriptions match primary ones. If we don't
// have a subscription, add it locally. If disabled on the Primary, disable it on the
// Secondary and log the activity. If enabled on the Primary but not on the secondary,
// it will do nothing. Finally, perform a data sync on all unlocked active slices.
func (s *Sync) syncSubscriptions(locals, prims subsArray, primaryID uuid.UUID) (err error) {
	// save our local subscriptions in a dictionary
	subs := make(sandpiper.SubsMap)
	subs.Load(locals)

	// run through primary (remote) subs and add to ours (local) if missing
	// (the slice for a new subscription is also added, but not the slice metadata,
	// which is added during the sync process)
	for _, remote := range prims {
		local, found := subs[remote.SubID]
		if !found {
			local = remote
			local.CompanyID = primaryID // change to our frame of reference for the add
			if err := s.sdb.AddSubscription(s.db, local); err != nil {
				return err
			}
			if err := s.sdb.AddSlice(s.db, remote.Slice); err != nil {
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
			if err := s.syncSlice(local.SubID, local.Slice, remote.Slice); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Sync) syncSlice(subID uuid.UUID, localSlice, remoteSlice *sandpiper.Slice) (err error) {
	// log activity a slice level *only* if an error occurs
	defer func(begin time.Time) {
		if err != nil {
			msg := "Slice \"" + localSlice.Name + "\""
			_ = s.sdb.LogActivity(s.db, subID, msg, time.Since(begin), err)
		}
	}(time.Now())

	if !remoteSlice.AllowSync {
		return errors.New("locked on server")
	}

	// todo: do we need to recalc local hash or can we just use the last one saved?
	if remoteSlice.ContentHash != localSlice.ContentHash {
		// update local metadata
		// get remote grain list
		// compare remote grain list with local grain list
		// for each new_oid in slice.oid_list
		//   get object.new_oid     // this object contains the payload
		//   store object in local.slice
		//   remove all obsolete local.slice.objects
	}

	return nil
}

// Subscriptions returns all subscriptions with slices and metadata (not paginated)
// for the current user's company
func (s *Sync) Subscriptions(c echo.Context) ([]sandpiper.Subscription, error) {
	if err := s.rbac.EnforceServerRole("primary"); err != nil {
		return nil, err
	}
	if err := s.rbac.EnforceRole(c, sandpiper.CompanyAdminRole); err != nil {
		return nil, err
	}
	companyID := s.rbac.CurrentUser(c).CompanyID
	return s.sdb.Subscriptions(s.db, companyID)
}

// Process (NOT CURRENTLY USED) responds to a sync start request and "upgrades" http to a websocket
func (s *Sync) Process(c echo.Context) error {
	if err := s.rbac.EnforceServerRole("primary"); err != nil {
		return err
	}
	if err := s.rbac.EnforceRole(c, sandpiper.CompanyAdminRole); err != nil {
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
