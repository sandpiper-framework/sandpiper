// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql

// tag service database access

import (
	"net/http"
	"strings"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/shared/model"
)

// Tag represents the client for tag table
type Tag struct{}

// NewTag returns a new tag database instance
func NewTag() *Tag {
	return &Tag{}
}

// Custom errors
var (
	// ErrAlreadyExists indicates the tag name is already used
	ErrAlreadyExists = echo.NewHTTPError(http.StatusInternalServerError, "Tag name already exists.")
)

// Create creates a new Tag in database (assumes allowed to do this)
func (s *Tag) Create(db orm.DB, tag sandpiper.Tag) (*sandpiper.Tag, error) {
	// don't add if duplicate name
	if nameExists(db, tag.Name) {
		return nil, ErrAlreadyExists
	}
	if err := db.Insert(&tag); err != nil {
		return nil, err
	}
	return &tag, nil
}

// View returns a single tag by ID (assumes allowed to do this)
func (s *Tag) View(db orm.DB, id int) (*sandpiper.Tag, error) {
	var tag = &sandpiper.Tag{ID: id}

	err := db.Select(tag)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

// List returns list of all tags
func (s *Tag) List(db orm.DB, p *sandpiper.Pagination) ([]sandpiper.Tag, error) {
	var subs []sandpiper.Tag

	err := db.Model(subs).Limit(p.Limit).Offset(p.Offset).Order("name").Select()
	if err != nil {
		return nil, err
	}
	return subs, nil
}

// Update updates tag info by primary key (assumes allowed to do this)
func (s *Tag) Update(db orm.DB, sub *sandpiper.Tag) error {
	_, err := db.Model(sub).UpdateNotZero()
	return err
}

// Delete removes the tag by primary key
func (s *Tag) Delete(db orm.DB, sub *sandpiper.Tag) error {
	return db.Delete(sub)
}

// nameExists returns true if name found in database
func nameExists(db orm.DB, name string) bool {
	m := new(sandpiper.Tag)
	err := db.Model(m).Where("lower(name) = ?", strings.ToLower(name)).Select()
	return err != pg.ErrNoRows
}
