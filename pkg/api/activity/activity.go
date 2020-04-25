// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package activity contains services for activity history.
package activity

import (
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/shared/model"
)

// Create makes a new sync activity record. Must be a sandpiper admin.
func (s *Activity) Create(c echo.Context, req sandpiper.Activity) (*sandpiper.Activity, error) {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return nil, err
	}
	return s.sdb.Create(s.db, req)
}

// View returns a single sync activity if allowed
func (s *Activity) View(c echo.Context, activityID int) (*sandpiper.Activity, error) {
	au := s.rbac.CurrentUser(c)
	if !au.AtLeast(sandpiper.AdminRole) {
		return nil, echo.ErrForbidden
	}
	return s.sdb.View(s.db, activityID)
}

// List returns list of sync activity scoped by user
func (s *Activity) List(c echo.Context, p *sandpiper.Pagination) ([]sandpiper.Activity, error) {
	// todo: should this only allow admin role (with no scoping)?
	q, err := s.rbac.EnforceScope(c)
	if err != nil {
		return nil, err
	}
	return s.sdb.List(s.db, q, p)
}

// Delete deletes a sync activity by id, if allowed
func (s *Activity) Delete(c echo.Context, id int) error {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return err
	}
	return s.sdb.Delete(s.db, id)
}
