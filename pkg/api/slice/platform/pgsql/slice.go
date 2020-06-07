// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package pgsql

// slice service database access.
// Manage slices and related metadata, but not which companies subscribe to the slice.

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
)

// Custom errors
var (
	ErrAlreadyExists = echo.NewHTTPError(http.StatusInternalServerError, "Slice name already exists.")
	ErrSliceNotFound = echo.NewHTTPError(http.StatusNotFound, "Slice not found.")
)

// Slice represents the client for slice table
type Slice struct{}

// NewSlice returns a new slice database instance
func NewSlice() *Slice {
	return &Slice{}
}

// sliceList holds multiple slice records returned from the database
type sliceList []sandpiper.Slice

// The IDs method creates an array of slice_ids
func (sl sliceList) IDs() []uuid.UUID {
	var ids = make([]uuid.UUID, 0, len(sl))

	for _, slice := range sl {
		ids = append(ids, slice.ID)
	}
	return ids
}

// FilterByTags creates a new Slices array using the tag query
func (sl sliceList) FilterByTags(db orm.DB, tags *sandpiper.TagQuery) (sliceList, error) {
	var tagged sliceList
	// See "Toxi" Solution http://howto.philippkeller.com/2005/04/24/Tags-Database-schemas/

	if tags.IsUnion { // slice must include at least one of the tags (union)
		err := db.Model(&tagged).Column("slice.id").
			Join("INNER JOIN slice_tags AS st ON slice.id = st.slice_id").
			Join("INNER JOIN tags AS t ON st.tag_id = t.id").
			Where("t.name IN (?)", pg.In(tags.TagList)).
			Group("slice.id").Select()
		if err != nil {
			return nil, err
		}
	} else { // slice must include all of the tags (intersection)
		err := db.Model(&tagged).Column("slice.id").
			Join("INNER JOIN slice_tags AS st ON slice.id = st.slice_id").
			Join("INNER JOIN tags AS t ON st.tag_id = t.id").
			Where("t.name IN (?)", pg.In(tags.TagList)).
			Group("slice.id").
			Having("COUNT(slice.id) = ?", tags.Count()).Select()
		if err != nil {
			return nil, err
		}
	}

	// put tagged slice ids in a set (map) for fast access
	taggedSet := make(map[uuid.UUID]bool)
	for _, slice := range tagged {
		taggedSet[slice.ID] = true
	}

	// run through received slice list "sl" and add to results if found in tagged set
	// uses "filtering without allocating" (https://github.com/golang/go/wiki/SliceTricks)
	results := sl[:0]
	for _, slice := range sl {
		if taggedSet[slice.ID] {
			results = append(results, slice)
		}
	}

	return results, nil
}

// Create adds a new slice with optional metadata (assumes allowed to do this)
func (s *Slice) Create(db orm.DB, slice sandpiper.Slice) (*sandpiper.Slice, error) {

	// don't add if would create a duplicate name
	if err := checkDuplicate(db, slice.Name); err != nil {
		return nil, err
	}

	// insert supplied slice data ("id" already assigned by create service)
	if err := db.Insert(&slice); err != nil {
		return nil, err
	}

	// insert any supplied meta data
	meta := &sandpiper.SliceMetadata{SliceID: slice.ID}
	for k, v := range slice.Metadata {
		meta.Key, meta.Value = k, v
		if err := db.Insert(meta); err != nil {
			return nil, err
		}
	}

	return &slice, nil
}

// View returns a single slice by ID with metadata and subscribed companies
func (s *Slice) View(db orm.DB, sliceID uuid.UUID) (*sandpiper.Slice, error) {
	// call view by subscription without a company
	return s.ViewBySub(db, uuid.Nil, sliceID)
}

// ViewBySub returns a single slice by ID if included in provided company subscriptions.
func (s *Slice) ViewBySub(db orm.DB, companyID uuid.UUID, sliceID uuid.UUID) (*sandpiper.Slice, error) {
	var slice = &sandpiper.Slice{ID: sliceID}

	// this filter function adds a condition to the "companies" relationship
	var filterFn = func(q *orm.Query) (*orm.Query, error) {
		if companyID == uuid.Nil {
			return q, nil
		}
		return q.Where("company_id = ?", companyID), nil
	}

	// get slice with subscribed companies
	err := db.Model(slice).Relation("Companies", filterFn).WherePK().Select()
	if err != nil {
		return nil, selectError(err)
	}

	// insert any metadata for the slice as a map
	slice.Metadata, err = metaDataMap(db, sliceID)

	return slice, err
}

// ViewByName returns a single slice by slice-name optionally limited by company subscriptions
func (s *Slice) ViewByName(db orm.DB, companyID uuid.UUID, name string) (*sandpiper.Slice, error) {
	var slice = &sandpiper.Slice{Name: name}

	// this filter function can add a condition to the "companies" relationship
	var filterFn = func(q *orm.Query) (*orm.Query, error) {
		if companyID == uuid.Nil {
			return q, nil
		}
		return q.Where("company_id = ?", companyID), nil
	}

	// get slice with subscribed companies
	err := db.Model(slice).
		Relation("Companies", filterFn).
		Where("lower(name) = ?", strings.ToLower(name)).
		Select()
	if err != nil {
		return nil, selectError(err)
	}

	// insert any metadata for the slice as a map
	slice.Metadata, err = metaDataMap(db, slice.ID)

	return slice, err
}

// List returns a list of all slices limited by scope and paginated
func (s *Slice) List(db orm.DB, tags *sandpiper.TagQuery, sc *sandpiper.Scope, p *sandpiper.Pagination) ([]sandpiper.Slice, error) {
	var slices sliceList

	// filter function adds optional condition to "Companies" relationship
	var filterFn = func(q *orm.Query) (*orm.Query, error) {
		if sc != nil {
			return q.Where(sc.Condition, sc.ID), nil
		}
		return q, nil
	}

	err := db.Model(&slices).Relation("Companies", filterFn).
		Limit(p.Limit).Offset(p.Offset).
		Order("name").Select()
	if err != nil {
		return nil, err
	}

	// create new Slices limited to TagQuery (if one was provided)
	if tags.Provided() {
		filtered, err := slices.FilterByTags(db, tags)
		if err != nil {
			return nil, err
		}
		slices = filtered
	}

	if len(slices) == 0 {
		return slices, nil
	}

	// look up metadata for all slices returned above (using an "in" list)
	var meta sandpiper.MetaArray
	ids := slices.IDs()

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

// Metadata returns an array of metadata for a slice
func (s *Slice) Metadata(db orm.DB, sliceID uuid.UUID) (sandpiper.MetaArray, error) {
	var meta sandpiper.MetaArray
	err := db.Model(&meta).Where("slice_id = ?", sliceID).Select()
	if err != nil {
		return nil, err
	}
	return meta, nil
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

// Refresh a slice's content information
func (s *Slice) Refresh(db orm.DB, sliceID uuid.UUID) error {

	// calculate hash of all grains in the slice
	hash, count, err := HashSlice(db, sliceID)
	if err != nil {
		return err
	}

	// update slice with new information
	slice := sandpiper.Slice{
		ID:           sliceID,
		ContentHash:  hash,
		ContentCount: count,
		ContentDate:  time.Now(),
	}
	_, err = db.Model(&slice).Column("content_hash", "content_count", "content_date").WherePK().Update()

	return err
}

// Lock disallows sync operations on this slice
func (s *Slice) Lock(db orm.DB, sliceID uuid.UUID) error {
	slice := sandpiper.Slice{ID: sliceID, AllowSync: false}
	_, err := db.Model(&slice).Column("allow_sync").WherePK().Update()
	return err
}

// Unlock allows sync operations on this slice
func (s *Slice) Unlock(db orm.DB, sliceID uuid.UUID) error {
	slice := sandpiper.Slice{ID: sliceID, AllowSync: true}
	_, err := db.Model(&slice).Column("allow_sync").WherePK().Update()
	return err
}

// metaDataMap returns a map of slice metadata. We use this separate query instead of
// an orm relationship because we don't want array of structs in json here.
// Maps marshal as {"key1": "value1", "key2": "value2", ...}
func metaDataMap(db orm.DB, sliceID uuid.UUID) (sandpiper.MetaMap, error) {
	var meta sandpiper.MetaArray
	err := db.Model(&meta).Where("slice_id = ?", sliceID).Select()
	if err != nil {
		return nil, err
	}
	return meta.ToMap(sliceID), nil
}

// checkDuplicate returns true if name found in database
func checkDuplicate(db orm.DB, name string) error {
	// attempt to select by unique key
	m := new(sandpiper.Slice)
	err := db.Model(m).
		Column("id").
		Where("lower(name) = ?", strings.ToLower(name)).
		Select()

	switch err {
	case pg.ErrNoRows: // ok to add
		return nil
	case nil: // found a row, so a duplicate
		return ErrAlreadyExists
	default: // return any other problem found
		return err
	}
}

func selectError(err error) error {
	if err == pg.ErrNoRows {
		return ErrSliceNotFound
	}
	return err
}

// HashSlice returns a sha1 hash of all metadata and grains in a slice
func HashSlice(db orm.DB, sliceID uuid.UUID) (string, int, error) {
	var (
		ids  []uuid.UUID
		b    bytes.Buffer
		meta sandpiper.MetaArray
	)

	var hash = func(meta sandpiper.MetaArray, ids []uuid.UUID) string {
		if len(meta) == 0 && len(ids) == 0 {
			return ""
		}
		for _, u := range ids {
			b.Write(u[:])
		}
		for _, m := range meta {
			b.WriteString(m.Key)
			b.WriteString(m.Value)
		}
		return fmt.Sprintf("%x", sha1.Sum(b.Bytes()))
	}

	// get slice metadata (sorted!)
	if err := db.Model(&meta).
		Where("slice_id = ?", sliceID).
		Order("key").
		Select(); err != nil {
		return "", 0, err
	}

	// get grain ids for the slice (sorted!)
	if err := db.Model().Table("grains").Column("id").
		Where("slice_id = ?", sliceID).
		Order("id").
		Select(&ids); err != nil {
		return "", 0, err
	}

	return hash(meta, ids), len(ids), nil
}
