// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package tag contains services for the tags resource.
package tag

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/shared/model"
)

// See "Toxi" solution in this article:
// http://howto.philippkeller.com/2005/04/24/Tags-Database-schemas/

// Create adds a new tag if administrator
func (s *Tag) Create(c echo.Context, req sandpiper.Tag) (*sandpiper.Tag, error) {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return nil, err
	}
	return s.sdb.Create(s.db, req)
}

// List returns list of tags that you can view
func (s *Tag) List(c echo.Context, p *sandpiper.Pagination) ([]sandpiper.Tag, error) {
	return s.sdb.List(s.db, p)
}

// View returns a single tag if allowed
func (s *Tag) View(c echo.Context, id int) (*sandpiper.Tag, error) {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return nil, err
	}
	return s.sdb.View(s.db, id)
}

// Delete deletes a tag if administrator
func (s *Tag) Delete(c echo.Context, id int) error {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return err
	}
	tag, err := s.sdb.View(s.db, id)
	if err != nil {
		return err
	}
	return s.sdb.Delete(s.db, tag)
}

// Update contains tag field request used for updating
type Update struct {
	ID          int
	Name        string
	Description string
}

// Update updates tag information
func (s *Tag) Update(c echo.Context, r *Update) (*sandpiper.Tag, error) {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return nil, err
	}

	tag := sandpiper.Tag{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}

	if err := s.sdb.Update(s.db, &tag); err != nil {
		return nil, err
	}
	return s.sdb.View(s.db, r.ID)
}

// Assign adds a tag assignment to a slice
func (s *Tag) Assign(c echo.Context, tagID int, sliceID uuid.UUID) error {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return err
	}
	return s.sdb.Assign(s.db, tagID, sliceID)
}

// Remove deletes a tag assignment from a slice
func (s *Tag) Remove(c echo.Context, tagID int, sliceID uuid.UUID) error {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return err
	}
	return s.sdb.Remove(s.db, tagID, sliceID)
}
