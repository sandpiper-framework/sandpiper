// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package grain contains services for grains
package grain

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/internal/model"
	"autocare.org/sandpiper/pkg/internal/scope"
)

// Create creates a new grain to hold data-objects
func (s *Grain) Create(c echo.Context, req sandpiper.Grain) (*sandpiper.Grain, error) {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return nil, err
	}
	return s.gdb.Create(s.db, req)
}

// List returns list of grains
func (s *Grain) List(c echo.Context, p *sandpiper.Pagination) ([]sandpiper.Grain, error) {
	au := s.rbac.CurrentUser(c)
	q, err := scope.Limit(au)
	if err != nil {
		return nil, err
	}
	return s.gdb.List(s.db, q, p)
}

// View returns a single grain if allowed
func (s *Grain) View(c echo.Context, grainID uuid.UUID) (*sandpiper.Grain, error) {
	au := s.rbac.CurrentUser(c)
	if au.Role != sandpiper.AdminRole {
		// make sure the grain is subscribed to this user's company
		return s.gdb.ViewBySub(s.db, au.CompanyID, grainID)
	}
	return s.gdb.View(s.db, grainID)
}

// Delete deletes a grain if allowed
func (s *Grain) Delete(c echo.Context, id uuid.UUID) error {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return err
	}
	grain, err := s.gdb.View(s.db, id)
	if err != nil {
		return err
	}
	return s.gdb.Delete(s.db, grain)
}
