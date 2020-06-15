// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package grain contains services for grains. Grains must belong to a slice
// and do not have an update method (use add/delete).
package grain

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/params"
)

// Create makes a new grain to hold our syncable data-objects. Must be a sandpiper admin.
func (s *Grain) Create(c echo.Context, replaceFlag bool, req *sandpiper.Grain) (*sandpiper.Grain, error) {
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

// ViewByKeys returns a single grain by (slice_id, grain_key) (admin function only)
func (s *Grain) ViewByKeys(c echo.Context, sliceID uuid.UUID, grainKey string, payloadFlag bool) (*sandpiper.Grain, error) {
	au := s.rbac.CurrentUser(c)
	if au.AtLeast(sandpiper.AdminRole) {
		return s.sdb.ViewByKeys(s.db, sliceID, grainKey, payloadFlag)
	}
	return nil, echo.ErrForbidden
}

// List returns list of grains scoped by user
func (s *Grain) List(c echo.Context, payload bool, p *params.Params) ([]sandpiper.Grain, error) {
	q, err := s.rbac.EnforceScope(c)
	if err != nil {
		return nil, err
	}
	return s.sdb.List(s.db, uuid.Nil, payload, q, p)
}

// ListBySlice returns a list of grains for a slice scoped by user
func (s *Grain) ListBySlice(c echo.Context, sliceID uuid.UUID, payloadFlag bool, p *params.Params) ([]sandpiper.Grain, error) {
	q, err := s.rbac.EnforceScope(c)
	if err != nil {
		return nil, err
	}
	return s.sdb.List(s.db, sliceID, payloadFlag, q, p)
}

// Delete deletes a grain by id, if allowed
func (s *Grain) Delete(c echo.Context, id uuid.UUID) error {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return err
	}
	return s.sdb.Delete(s.db, id)
}
