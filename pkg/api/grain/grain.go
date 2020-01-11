// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package grain contains services for grains. Grains must belong to a slice
// and do not have an update method (use add/delete).
package grain

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/internal/model"
	"autocare.org/sandpiper/pkg/internal/scope"
)

// Create creates a new grain to hold data-objects (we pass request as a pointer
// for this service because it can be very big).
func (s *Grain) Create(c echo.Context, req *sandpiper.Grain) (*sandpiper.Grain, error) {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return nil, err
	}
	return s.sdb.Create(s.db, req)
}

// List returns list of grains scoped by user
func (s *Grain) List(c echo.Context, p *sandpiper.Pagination) ([]sandpiper.Grain, error) {
	au := s.rbac.CurrentUser(c)
	q, err := scope.Limit(au)
	if err != nil {
		return nil, err
	}
	return s.sdb.List(s.db, q, p)
}

// View returns a single grain if allowed
func (s *Grain) View(c echo.Context, grainID uuid.UUID) (*sandpiper.Grain, error) {
	au := s.rbac.CurrentUser(c)
	if au.Role != sandpiper.AdminRole {
		// make sure the grain is subscribed to this user's company
		return s.sdb.ViewBySub(s.db, au.CompanyID, grainID)
	}
	return s.sdb.View(s.db, grainID)
}

// Delete deletes a grain if allowed
func (s *Grain) Delete(c echo.Context, id uuid.UUID) error {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return err
	}
	grain, err := s.sdb.View(s.db, id)
	if err != nil {
		return err
	}
	return s.sdb.Delete(s.db, grain)
}
