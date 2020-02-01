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
	ErrAlreadyExists     = echo.NewHTTPError(http.StatusInternalServerError, "Tag name already exists.")
	ErrTagDoesNotExist   = echo.NewHTTPError(http.StatusInternalServerError, "Cannot assign a tag ID that does not exist.")
	ErrSliceDoesNotExist = echo.NewHTTPError(http.StatusInternalServerError, "Cannot assign to a slice ID that does not exist.")
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

// Assign adds a a slice_tag record
func (s *Tag) Assign(db orm.DB, tagID int, sliceID uuid.UUID) error {
	sliceTag := sandpiper.SliceTag{TagID: tagID, SliceID: sliceID}
	if err := db.Insert(&sliceTag); err != nil {
		if chk := checkAssignIDs(db, tagID, sliceID); chk != nil {
			return chk
		}
		return err
	}
	return nil
}

// Remove deletes a slice_tag record by composite primary key
func (s *Tag) Remove(db orm.DB, tagID int, sliceID uuid.UUID) error {
	sliceTag := sandpiper.SliceTag{TagID: tagID, SliceID: sliceID}
	return db.Delete(sliceTag)
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

func checkAssignIDs(db orm.DB, tagID int, sliceID uuid.UUID) error {
	if !tagExists(db, tagID) {
		return ErrTagDoesNotExist
	}
	if !sliceExists(db, sliceID) {
		return ErrSliceDoesNotExist
	}
	return nil
}

func tagExists(db orm.DB, tagID int) bool {
	m := &sandpiper.Tag{ID: tagID}
	err := db.Model(m).Column("id").WherePK().Select()
	// return true if found (nil) or any error besides not found
	// return false if not found (ErrNoRows) or any other error
	return (err == nil) || (err != pg.ErrNoRows)
}

func sliceExists(db orm.DB, sliceID uuid.UUID) bool {
	m := &sandpiper.Slice{ID: sliceID}
	err := db.Model(m).Column("id").WherePK().Select()
	// return true if found (nil) or any error besides not found
	// return false if not found (ErrNoRows) or any other error
	return (err == nil) || (err != pg.ErrNoRows)
}
