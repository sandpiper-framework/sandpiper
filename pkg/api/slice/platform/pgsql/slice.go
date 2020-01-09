// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql

// slice service database access

import (
	"net/http"
	"strings"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/internal/model"
	"autocare.org/sandpiper/pkg/internal/scope"
)

// Slice represents the client for slice table
type Slice struct{}

// NewSlice returns a new slice database instance
func NewSlice() *Slice {
	return &Slice{}
}

// Custom errors
var (
	// ErrAlreadyExists indicates the slice name is already used
	ErrAlreadyExists = echo.NewHTTPError(http.StatusInternalServerError, "Slice name already exists.")
)

// Create creates a new slice in database (assumes allowed to do this)
func (s *Slice) Create(db orm.DB, slice sandpiper.Slice) (*sandpiper.Slice, error) {

	// don't add if would create a duplicate name
	if nameExists(db, slice.Name) {
		return nil, ErrAlreadyExists
	}

	// insert supplied slice data ("id" already assigned by create service)
	if err := db.Insert(&slice); err != nil {
		return nil, err
	}

	// insert any supplied meta data
	meta := &sandpiper.SliceMetadata{SliceID: slice.ID}
	for k, v := range slice.Metadata {
		meta.Key = k
		meta.Value = v
		if err := db.Insert(meta); err != nil {
			return nil, err
		}
	}

	return &slice, nil
}

// View returns a single slice by ID with metadata and subscribed companies
// (optionally limited to a company)
func (s *Slice) View(db orm.DB, sliceID uuid.UUID) (*sandpiper.Slice, error) {
	var slice = &sandpiper.Slice{ID: sliceID}

	// get slice by primary key with subscribed companies
	err := db.Model(slice).Relation("Companies").WherePK().Select()
	if err != nil {
		return nil, err
	}

	// insert any metadata via a map
	slice.Metadata, err = metaDataMap(db, slice)

	return slice, nil
}

// ViewBySub returns a single slice by ID if included in provided company subscriptions.
func (s *Slice) ViewBySub(db orm.DB, companyID uuid.UUID, sliceID uuid.UUID) (*sandpiper.Slice, error) {
	var slice = &sandpiper.Slice{ID: sliceID}

	// this filter function adds a condition to the companies relationship
	var filterFn = func(q *orm.Query) (*orm.Query, error) {
		return q.Where("company_id = ?", companyID), nil
	}

	// get slice with subscribed companies
	err := db.Model(slice).Column("slice.*").Relation("Companies", filterFn).WherePK().Select()
	if err != nil {
		return nil, err
	}

	// insert any metadata via a map
	slice.Metadata, err = metaDataMap(db, slice)

	return slice, err
}

// List returns a list of all slices limited by scope and paginated
func (s *Slice) List(db orm.DB, sc *scope.Clause, p *sandpiper.Pagination) ([]sandpiper.Slice, error) {
	var slices sandpiper.SliceArray

	// this filter function adds an optional condition to the companies relationship
	var filterFn = func(q *orm.Query) (*orm.Query, error) {
		if sc != nil {
			return q.Where(sc.Condition, sc.ID), nil
		}
		return q, nil
	}

	err := db.Model(&slices).Relation("Companies", filterFn).
		Limit(p.Limit).Offset(p.Offset).Order("name").Select()
	if err != nil {
		return nil, err
	}

	// look up metadata for all slices returned above (using an "in" list)
	var meta sandpiper.MetaArray
	ids := slices.SliceIDs()

	err = db.Model(&meta).Where("slice_id in (?)", pg.In(ids)).Select()
	if err != nil {
		return nil, err
	}

	// insert metadata for each slice into response
	for i, slice := range slices {
		slices[i].Metadata = meta.ToMap(slice.ID)
	}

	return slices, nil
}

// Update updates slice info by primary key (assumes allowed to do this)
func (s *Slice) Update(db orm.DB, slice *sandpiper.Slice) error {
	// todo: should we delete all metadata and re-add from the map?
	_, err := db.Model(slice).Update()
	return err
}

// Delete a slice
func (s *Slice) Delete(db orm.DB, slice *sandpiper.Slice) error {
	// WARNING: Foreign key constraints remove related metadata and grains!
	return db.Delete(slice)
}

// metaDataMap returns a map of slice metadata. We use this separate query instead of
// an orm relationship because we don't want array of structs in json. Maps will marshal
// as {"key1": "value1", "key2": "value2", ...}
func metaDataMap(db orm.DB, slice *sandpiper.Slice) (sandpiper.MetaMap, error) {
	var meta sandpiper.MetaArray
	err := db.Model(&meta).Where("slice_id = ?", slice.ID).Select()
	if err != nil {
		return nil, err
	}
	return meta.ToMap(slice.ID), nil
}

// nameExists returns true if name found in database
func nameExists(db orm.DB, name string) bool {
	m := new(sandpiper.Slice)
	err := db.Model(m).
		Column("id", "name").
		Where("lower(name) = ?", strings.ToLower(name)).
		Select()
	return err != pg.ErrNoRows
}
