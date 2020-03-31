// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package grain contains services for grains. Grains must belong to a slice
// and do not have an update method (use add/delete).
package grain

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/shared/model"
)

// Create makes a new grain to hold our syncable data-objects. Must be a sandpiper admin.
func (s *Grain) Create(c echo.Context, replaceFlag bool, req sandpiper.Grain) (*sandpiper.Grain, error) {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return nil, err
	}
	return s.sdb.Create(s.db, replaceFlag, req)
}

// View returns a single grain if allowed
func (s *Grain) View(c echo.Context, grainID uuid.UUID) (*sandpiper.Grain, error) {
	au := s.rbac.CurrentUser(c)
	if !au.AtLeast(sandpiper.AdminRole) {
		// make sure the grain is subscribed to this user's company
		if !s.sdb.CompanySubscribed(s.db, au.CompanyID, grainID) {
			return nil, echo.ErrForbidden
		}
	}
	return s.sdb.View(s.db, grainID)
}

// Exists returns a single grain's basic information (without authorization checks)
func (s *Grain) Exists(c echo.Context, sliceID uuid.UUID, grainKey string) (*sandpiper.Grain, error) {
	return s.sdb.Exists(s.db, sliceID, grainKey)
}

// List returns list of grains scoped by user
func (s *Grain) List(c echo.Context, payload bool, p *sandpiper.Pagination) ([]sandpiper.Grain, error) {
	q, err := s.rbac.EnforceScope(c)
	if err != nil {
		return nil, err
	}
	return s.sdb.List(s.db, payload, q, p)
}

// Delete deletes a grain by id, if allowed
func (s *Grain) Delete(c echo.Context, id uuid.UUID) error {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return err
	}
	return s.sdb.Delete(s.db, id)
}
