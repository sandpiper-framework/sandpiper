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
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/sync/credentials"
	"autocare.org/sandpiper/pkg/shared/client"
	"autocare.org/sandpiper/pkg/shared/model"
)

// Start sends a sync request to a primary sandpiper server from our secondary server
func (s *Sync) Start(c echo.Context, primaryID uuid.UUID) error {
	// must be a secondary server to start the sync
	if err := s.rbac.EnforceServerRole("secondary"); err != nil {
		return err
	}
	// must be a local admin to start the sync
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return err
	}

	// get company information for the primary server we're syncing
	p, err := s.sdb.Primary(s.db, primaryID)
	if err != nil {
		return err
	}

	// get login credentials from company's sync_api_key
	creds, err := credentials.New(p.SyncAPIKey, s.sec.APIKeySecret())
	if err != nil {
		return err
	}

	// login to the primary server (saving token)
	addr, err := url.ParseRequestURI(p.SyncAddr)
	if err != nil {
		return err
	}
	api, err := client.Login(addr, creds.User, creds.Password, false)
	if err != nil {
		return err
	}

	// upgrade to websocket (as a client) and perform the sync process
	if err := api.Process(); err != nil {
		return err
	}

	// ask for our subscriptions

	// add any new subscriptions

	/*
		for _, sub := range subs {
			slice := sub.Slice
			if !slice.AllowSync {
				// just log that it is locked
			}

		}
	*/

	// send a "GET /sync" request to the primary server address
	// ask for subscriptions (add any not already in our database -- with empty contents)
	/*
		for each slice in my.slices
			get slice.hash    // sha1 hash representing slice oids
			if slice.hash <> local.slice.hash    // saved from last sync or must we recalc?
					get slice.oid_list
			compare slice.oid_list with local.slice.oid_list
			for each new_oid in slice.oid_list
				get object.new_oid     // this object contains the payload
				store object in local.slice
			remove all obsolete local.slice.objects
		put slice.sync_completed   // for activity reporting purposes
	*/

	return nil
}

// Process responds to a sync start request and "upgrades" http to a websocket
func (s *Sync) Process(c echo.Context) error {
	var (
		upgrader = websocket.Upgrader{}
	)

	if err := s.rbac.EnforceServerRole("primary"); err != nil {
		return err
	}
	if err := s.rbac.EnforceRole(c, sandpiper.CompanyAdminRole); err != nil {
		return err
	}
	// todo: do the actual work here, making calls to s.sdb.xxx as necessary

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

	return nil
}
