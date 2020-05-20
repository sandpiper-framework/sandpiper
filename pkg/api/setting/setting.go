// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package setting contains services for settings. Settings must belong to a slice
// and do not have an update method (use add/delete).
package setting

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"sandpiper/pkg/shared/model"
)

// Create makes a new setting to hold our syncable data-objects. Must be a sandpiper admin.
func (s *Setting) Create(c echo.Context, req *sandpiper.Setting) (*sandpiper.Setting, error) {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return nil, err
	}
	return s.sdb.Create(s.db, req)
}

// View returns the database settings record
func (s *Setting) View(c echo.Context) (*sandpiper.Setting, error) {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return nil, err
	}
	return s.sdb.View(s.db)
}

// Update contains settings request used for updating
type Update struct {
	ID         bool
	ServerRole string
	ServerID   uuid.UUID
}

// Update updates setting information
func (s *Setting) Update(c echo.Context, r *Update) (*sandpiper.Setting, error) {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return nil, err
	}
	set := sandpiper.Setting{
		ID:         true,
		ServerRole: r.ServerRole,
		ServerID:   r.ServerID,
	}
	err := s.sdb.Update(s.db, &set)
	if err != nil {
		return nil, err
	}
	return s.sdb.View(s.db)
}