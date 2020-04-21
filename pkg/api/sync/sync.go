// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package sync contains services for syncs. Syncs must belong to a slice
// and do not have an update method (use add/delete).
package sync

import (
	"github.com/labstack/echo/v4"
	"net/url"

	"autocare.org/sandpiper/pkg/shared/model"
)

// Start sends a sync request to a primary sandpiper server
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
