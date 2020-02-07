// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package slice contains services for slices
package slice

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"

	"autocare.org/sandpiper/pkg/shared/model"
)

// Custom errors
var (
	// ErrTagsNotAllowed indicates the slice name is already used
	ErrTagsNotAllowed = echo.NewHTTPError(http.StatusInternalServerError, "Not authorized for tagged queries")
)

// Create creates a new slice to hold data-objects
func (s *Slice) Create(c echo.Context, req sandpiper.Slice) (*sandpiper.Slice, error) {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return nil, err
	}
	return s.sdb.Create(s.db, req)
}

// List returns list of slices
func (s *Slice) List(c echo.Context, tags *sandpiper.TagQuery, p *sandpiper.Pagination) ([]sandpiper.Slice, error) {
	q, err := s.rbac.EnforceScope(c)
	if err != nil {
		return nil, err
	}
	au := s.rbac.CurrentUser(c)
	if tags.Provided() && !au.AtLeast(sandpiper.AdminRole) {
		// only admin can query by tags (internal structuring of slices)
		return nil, ErrTagsNotAllowed
	}
	return s.sdb.List(s.db, tags, q, p)
}

// View returns a single slice if allowed
func (s *Slice) View(c echo.Context, sliceID uuid.UUID) (*sandpiper.Slice, error) {
	au := s.rbac.CurrentUser(c)
	if au.AtLeast(sandpiper.AdminRole) {
		return s.sdb.View(s.db, sliceID)
	}
	// make sure the slice is subscribed to the current user's company
	return s.sdb.ViewBySub(s.db, au.CompanyID, sliceID)
}

// Update contains slice's information used for updating
type Update struct {
	ID           uuid.UUID
	Name         string
	ContentHash  string
	ContentCount uint
	ContentDate  time.Time
}

// Update updates slice information
func (s *Slice) Update(c echo.Context, r *Update) (*sandpiper.Slice, error) {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return nil, err
	}
	slice := &sandpiper.Slice{
		ID:           r.ID,
		Name:         r.Name,
		ContentHash:  r.ContentHash,
		ContentCount: r.ContentCount,
		ContentDate:  r.ContentDate,
	}
	err := s.sdb.Update(s.db, slice)
	if err != nil {
		return nil, err
	}
	return s.sdb.View(s.db, r.ID)
}

// Delete deletes a slice if allowed
func (s *Slice) Delete(c echo.Context, id uuid.UUID) error {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return err
	}
	slice, err := s.sdb.View(s.db, id)
	if err != nil {
		return err
	}
	return s.sdb.Delete(s.db, slice)
}
