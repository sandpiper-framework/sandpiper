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
	"net/url"

	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/shared/model"
)

// Start sends a sync request to a primary sandpiper server from our secondary server
func (s *Sync) Start(c echo.Context, primary *url.URL) error {
	// must be a secondary server to start the sync
	if err := s.rbac.EnforceServerRole("secondary"); err != nil {
		return err
	}
	// must be a local admin to start the sync
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return err
	}

	// todo: do the actual work here, making calls to s.sdb.xxx as necessary

	// login as a client (/login)

	// open websocket (as a client) with the remote primary server (/sync/{url})

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

// Connect responds to a sync start request and "upgrades" http to a websocket
func (s *Sync) Connect(c echo.Context) error {
	if err := s.rbac.EnforceServerRole("primary"); err != nil {
		return err
	}
	if err := s.rbac.EnforceRole(c, sandpiper.CompanyAdminRole); err != nil {
		return err
	}
	// todo: do the actual work here, making calls to s.sdb.xxx as necessary

	return nil
}
