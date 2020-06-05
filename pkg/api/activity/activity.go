// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package activity contains services for activity history.
package activity

import (
	"github.com/labstack/echo/v4"

	"sandpiper/pkg/shared/model"
)

// Create makes a new sync activity record (most often by a secondary sync user).
func (s *Activity) Create(c echo.Context, req sandpiper.Activity) (*sandpiper.Activity, error) {
	if err := s.rbac.EnforceRole(c, sandpiper.SyncRole); err != nil {
		return nil, err
	}
	return s.sdb.Create(s.db, req)
}

// View returns a single sync activity if allowed
func (s *Activity) View(c echo.Context, activityID int) (*sandpiper.Activity, error) {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return nil, err
	}
	return s.sdb.View(s.db, activityID)
}

// List returns list of sync activity scoped by user
func (s *Activity) List(c echo.Context, p *sandpiper.Pagination) ([]sandpiper.Activity, error) {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return nil, err
	}
	return s.sdb.List(s.db, p)
}

// Delete deletes a sync activity by id, if allowed
func (s *Activity) Delete(c echo.Context, id int) error {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return err
	}
	return s.sdb.Delete(s.db, id)
}
