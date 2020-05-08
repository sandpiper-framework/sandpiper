// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql

// sync service database access

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"sandpiper/pkg/shared/model"
)

// Custom errors
var (
	ErrAlreadyExists = echo.NewHTTPError(http.StatusInternalServerError, "Subscription name already exists.")
	ErrNoAccess      = echo.NewHTTPError(http.StatusForbidden, "No Access to this slice")
)

// Sync represents the client for sync table
type Sync struct{}

// NewSync returns a new sync instance
func NewSync() *Sync {
	return &Sync{}
}

// LogActivity adds a sync log entry to the activity table
func (s *Sync) LogActivity(db orm.DB, subID uuid.UUID, msg string, d time.Duration, err error) error {
	if err != nil {
		// append error to message
		msg = fmt.Sprintf("%s (%v)", msg, err)
	}
	// log this event in our activity table
	activity := sandpiper.Activity{
		SubID:    subID,
		Success:  err == nil,
		Message:  msg,
		Duration: d,
	}
	if err := db.Insert(&activity); err != nil {
		return err
	}
	return nil
}

// Primary returns a single (primary) company by ID (assumes allowed to do this)
func (s *Sync) Primary(db orm.DB, id uuid.UUID) (*sandpiper.Company, error) {
	var company = &sandpiper.Company{ID: id}

	err := db.Model(company).WherePK().Select()
	if err != nil {
		return nil, err
	}
	return company, nil
}

// Subscriptions returns list of all local subscriptions (with slice & metadata) for a company
// without pagination
func (s *Sync) Subscriptions(db orm.DB, companyID uuid.UUID) ([]sandpiper.Subscription, error) {
	var subs []sandpiper.Subscription

	err := db.Model(&subs).Relation("Slice").Where("company_id = ?", companyID).Select()
	if err != nil {
		return nil, err
	}
	return subs, nil
}

// AddSubscription creates a new Subscription in database
func (s *Sync) AddSubscription(db orm.DB, sub sandpiper.Subscription) error {
	// make sure name is unique on our side too
	if err := checkDupSubName(db, sub.Name); err != nil {
		if err != ErrAlreadyExists {
			return err
		}
		sub.Name = sub.Name + " (" + sub.SubID.String() + ")"
	}
	if err := db.Insert(&sub); err != nil {
		return err
	}
	return nil
}

// DeactivateSubscription turns off the active flag and logs the event in our activity
func (s *Sync) DeactivateSubscription(db orm.DB, subID uuid.UUID) error {
	sub := &sandpiper.Subscription{SubID: subID, Active: false}
	if _, err := db.Model(sub).Column("active").WherePK().Update(); err != nil {
		return err
	}
	// log this event in our activity table
	activity := sandpiper.Activity{
		SubID:   subID,
		Success: false,
		Message: "deactivated by primary",
	}
	if err := db.Insert(&activity); err != nil {
		return err
	}
	return nil
}

// AddSlice creates a new Slice in the database (without metadata)
func (s *Sync) AddSlice(db orm.DB, slice *sandpiper.Slice) error {
	// make sure name is unique on our side too
	if err := checkDupSliceName(db, slice.Name); err != nil {
		if err != ErrAlreadyExists {
			return err
		}
		slice.Name = slice.Name + " (" + slice.ID.String() + ")"
	}
	if err := db.Insert(&slice); err != nil {
		return err
	}
	return nil
}

// Slice returns a slice and its metadata
func (s *Sync) Slice(db orm.DB, sliceID uuid.UUID) (*sandpiper.Slice, error) {
	slice := &sandpiper.Slice{ID: sliceID}
	err := db.Model(slice).WherePK().Select()
	if err != nil {
		return nil, err
	}
	// insert any metadata for the slice as a map
	slice.Metadata, err = metaDataMap(db, slice.ID)

	return slice, nil
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

// SliceAccess checks if a slice is included in a company's subscriptions.
func (s *Sync) SliceAccess(db orm.DB, companyID uuid.UUID, sliceID uuid.UUID) error {
	sub := new(sandpiper.Subscription)
	err := db.Model(sub).Column("sub_id").
		Where("slice_id = ?", sliceID).
		Where("company_id = ?", companyID).
		Select()
	switch err {
	case pg.ErrNoRows:
		return ErrNoAccess
	case nil: // found a row, so have access
		return nil
	default: // return any other problem found
		return err
	}
}

// UpdateSliceMetadata replaces target metadata with the source metadata
func (s *Sync) UpdateSliceMetadata(db orm.DB, target, source *sandpiper.Slice) error {
	// see if there are any changes to make
	if target.Metadata.Equals(source.Metadata) {
		return nil
	}
	// drop existing target slice metadata
	md := new(sandpiper.SliceMetadata)
	if _, err := db.Model(md).Where("slice_id = ?", target.ID).Delete(); err != nil {
		return err
	}
	// add new source slice metadata
	meta := &sandpiper.SliceMetadata{SliceID: source.ID}
	for k, v := range source.Metadata {
		meta.Key, meta.Value = k, v
		if err := db.Insert(meta); err != nil {
			return err
		}
	}
	target.Metadata = source.Metadata
	return nil
}

// Grains returns a list of grains for a slice (with brief or all fields)
// assumes allowed to do this
func (s *Sync) Grains(db orm.DB, companyID uuid.UUID, sliceID uuid.UUID, briefFlag bool) ([]sandpiper.Grain, error) {
	var grains []sandpiper.Grain

	// columns to select
	cols := "grain.id"
	if !briefFlag {
		cols = cols + ", slice_id, grain_key, source, encoding, grain.created_at, length(payload) AS payload_len"
	}
	err := db.Model(&grains).Where("slice_id = ?", sliceID).Select()
	if err != nil {
		return nil, err
	}
	return grains, nil
}

// AddGrain adds a grain locally
func (s *Sync) AddGrain(db orm.DB, grain *sandpiper.Grain) error {
	if err := db.Insert(grain); err != nil {
		return err
	}
	return nil
}

// DeleteGrains removes all provided grain ids
func (s *Sync) DeleteGrains(db orm.DB, ids []uuid.UUID) error {
	_, err := db.Model((*sandpiper.Grain)(nil)).Where("id in (?)", pg.In(ids)).Delete()
	return err
}

// checkDupSubName returns true if name found in database
func checkDupSubName(db orm.DB, name string) error {
	// attempt to select by unique key
	m := new(sandpiper.Subscription)
	err := db.Model(m).
		Column("sub_id").
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

// checkDupSliceName returns true if name found in database
func checkDupSliceName(db orm.DB, name string) error {
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
