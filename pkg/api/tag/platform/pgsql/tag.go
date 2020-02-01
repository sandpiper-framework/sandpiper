// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql

// tag service database access

import (
	"net/http"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
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
	if err := checkDuplicate(db, tag.Name); err != nil {
		return nil, err
	}
	if err := db.Insert(&tag); err != nil {
		return nil, err
	}
	return &tag, nil
}

// Assign adds a tag to a slice
func (s *Tag) Assign(db orm.DB, tagID int, sliceID uuid.UUID) (*sandpiper.SliceTag, error) {
	sliceTag := sandpiper.SliceTag{TagID: tagID, SliceID: sliceID}
	if err := db.Insert(&sliceTag); err != nil {
		return nil, err
	}
	return &sliceTag, nil
}

// View returns a single tag by ID with any associated slices (assumes allowed to do this)
func (s *Tag) View(db orm.DB, id int) (*sandpiper.Tag, error) {
	var tag = &sandpiper.Tag{ID: id}

	err := db.Model(tag).Relation("Slices").WherePK().Select()
	if err != nil {
		return nil, err
	}
	return tag, nil
}

// List returns list of all tags
func (s *Tag) List(db orm.DB, p *sandpiper.Pagination) ([]sandpiper.Tag, error) {
	var tags []sandpiper.Tag

	err := db.Model(&tags).Limit(p.Limit).Offset(p.Offset).Order("name").Select()
	if err != nil {
		return nil, err
	}
	return tags, nil
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

// checkDuplicate returns true if name found in database
func checkDuplicate(db orm.DB, name string) error {
	// attempt to select by unique key
	m := new(sandpiper.Tag)
	err := db.Model(m).
		Column("id").
		Where("name = ?", name). // already lower-case
		Select()

	switch err {
	case pg.ErrNoRows: // ok to add
		return nil
	case nil: // would be a duplicate
		return ErrAlreadyExists
	default: // return any other problem found
		return err
	}
}
