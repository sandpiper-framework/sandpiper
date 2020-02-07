// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql

// slice service database access.
// Manage slices and related metadata, but not which companies subscribe to the slice.

import (
	"net/http"
	"strings"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/shared/model"
)

// Custom errors
var (
	// ErrAlreadyExists indicates the slice name is already used
	ErrAlreadyExists = echo.NewHTTPError(http.StatusInternalServerError, "Slice name already exists.")
)

// Slice represents the client for slice table
type Slice struct{}

// NewSlice returns a new slice database instance
func NewSlice() *Slice {
	return &Slice{}
}

// Slices holds multiple slice records returned from the database
type Slices []sandpiper.Slice

// The IDs method creates an array of slice_ids
func (a Slices) IDs() []uuid.UUID {
	var ids = make([]uuid.UUID, 0, len(a))

	for _, slice := range a {
		ids = append(ids, slice.ID)
	}
	return ids
}

// FilterByTags creates a new Slices array using the tag query
func (a Slices) FilterByTags(db orm.DB, tags *sandpiper.TagQuery) (Slices, error) {
	var tagged Slices

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

	// run through all received slices "a" and if found in tagged set, add to results
	// use "filtering without allocating" (https://github.com/golang/go/wiki/SliceTricks)
	results := a[:0]
	for _, slice := range a {
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
// (optionally limited to a company)
func (s *Slice) View(db orm.DB, sliceID uuid.UUID) (*sandpiper.Slice, error) {
	var slice = &sandpiper.Slice{ID: sliceID}

	// get slice by primary key with subscribed companies
	err := db.Model(slice).Relation("Companies").WherePK().Select()
	if err != nil {
		return nil, err
	}

	// insert any metadata for the slice as a map
	slice.Metadata, err = metaDataMap(db, sliceID)

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

	// insert any metadata for the slice as a map
	slice.Metadata, err = metaDataMap(db, sliceID)

	return slice, err
}

/*
todo: use these queries as a basis for filtering slices by tags (maybe as a CTE?)

CTE

WITH "scope" AS (
	SELECT slices.id FROM Slices
	INNER JOIN slice_tags ON Slices.id = slice_tags.slice_id
	INNER JOIN tags ON slice_tags.tag_id = tags.id
	WHERE tags.name IN ('brake_products', 'wiper_products')
	GROUP By slices.id
)
SELECT "slice"."id", "slice"."name", "slice"."content_hash", "slice"."content_count", "slice"."content_date", "slice"."created_at", "slice"."updated_at"
FROM "scope" JOIN "slices" AS "slice" ON slice.id = scope.id ORDER BY "name" LIMIT 100

CURRENT:

SELECT "slice"."id", "slice"."name", "slice"."content_hash", "slice"."content_count", "slice"."content_date", "slice"."created_at", "slice"."updated_at"
FROM "slices" AS "slice" ORDER BY "name" LIMIT 100

SELECT "subscriptions".*, "company"."id", "company"."name", "company"."sync_addr", "company"."active", "company"."created_at", "company"."updated_at"
FROM "companies" AS "company"
JOIN "subscriptions" AS "subscriptions" ON ("subscriptions"."slice_id") IN ('1b40204a-7acd-4c78-a3c4-0fa95d2f00f6', '2bea8308-1840-4802-ad38-72b53e31594c')
WHERE ("company"."id" = "subscriptions"."company_id")

SELECT "slice_metadata"."slice_id", "slice_metadata"."key", "slice_metadata"."value"
FROM "slice_metadata" AS "slice_metadata"
WHERE (slice_id in ('1b40204a-7acd-4c78-a3c4-0fa95d2f00f6','2bea8308-1840-4802-ad38-72b53e31594c'))

*/

// List returns a list of all slices limited by scope and paginated
func (s *Slice) List(db orm.DB, tags *sandpiper.TagQuery, sc *sandpiper.Scope, p *sandpiper.Pagination) ([]sandpiper.Slice, error) {
	var slices Slices

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
