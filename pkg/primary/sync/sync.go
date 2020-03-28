// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package sync contains services for sync history.
package sync

import (
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/shared/model"
)

// Create makes a new sync to hold our syncable data-objects. Must be a sandpiper admin.
func (s *Sync) Create(c echo.Context, req sandpiper.Sync) (*sandpiper.Sync, error) {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return nil, err
	}
	return s.sdb.Create(s.db, req)
}

// View returns a single sync if allowed
func (s *Sync) View(c echo.Context, syncID int) (*sandpiper.Sync, error) {
	au := s.rbac.CurrentUser(c)
	if !au.AtLeast(sandpiper.AdminRole) {
		// make sure the sync is subscribed to this user's company
		//if !s.sdb.CompanySubscribed(s.db, au.CompanyID, syncID) {
		//	return nil, echo.ErrForbidden
		//}
	}
	return s.sdb.View(s.db, syncID)
}

// List returns list of syncs scoped by user
func (s *Sync) List(c echo.Context, p *sandpiper.Pagination) ([]sandpiper.Sync, error) {
	q, err := s.rbac.EnforceScope(c)
	if err != nil {
		return nil, err
	}
	return s.sdb.List(s.db, q, p)
}

// Delete deletes a sync by id, if allowed
func (s *Sync) Delete(c echo.Context, id int) error {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return err
	}
	return s.sdb.Delete(s.db, id)
}
