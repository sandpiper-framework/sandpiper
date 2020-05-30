// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql

// sync service database access

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	slicesvc "sandpiper/pkg/api/slice/platform/pgsql"
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
	var errMsg string
	if err != nil {
		errMsg = fmt.Sprintf("%v", err)
	}
	activity := sandpiper.Activity{
		SubID:    subID,
		Success:  err == nil,
		Message:  msg,
		Error:    errMsg,
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

// Subscriptions returns list of all local subscriptions (with slice but not metadata) for a company
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
	if err := db.Insert(slice); err != nil {
		return err
	}
	return nil
}

// RefreshSlice updates the content fields and checks against source slice to make
// sure the sync agrees
func (s *Sync) RefreshSlice(db orm.DB, slice *sandpiper.Slice) error {

	// use a public function from the slice service to get our slice information
	hash, count, err := slicesvc.HashSliceGrains(db, slice.ID)
	if err != nil {
		return err
	}

	// update slice with new information
	m := sandpiper.Slice{
		ID:           slice.ID,
		ContentHash:  hash,
		ContentCount: count,
		ContentDate:  slice.ContentDate, // keep source's date
	}
	_, err = db.Model(&m).Column("content_hash", "content_count", "content_date").WherePK().Update()
	if err != nil {
		return err
	}

	// see if the sync worked (hash values match, etc.)
	if slice.ContentHash != hash || slice.ContentCount != count {
		return errors.New("content hash or count do not match after sync")
	}

	return err
}

// SliceMetadata returns slice metadata for a sliceID
func (s *Sync) SliceMetadata(db orm.DB, sliceID uuid.UUID) (sandpiper.MetaArray, error) {
	var meta sandpiper.MetaArray
	err := db.Model(&meta).Where("slice_id = ?", sliceID).Select()
	if err != nil {
		return nil, err
	}
	return meta, nil
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

// ReplaceSliceMetadata replaces target metadata with the source metadata
func (s *Sync) ReplaceSliceMetadata(db orm.DB, sliceID uuid.UUID, metaArray sandpiper.MetaArray) error {
	// drop existing slice metadata
	md := new(sandpiper.SliceMetadata)
	if _, err := db.Model(md).Where("slice_id = ?", sliceID).Delete(); err != nil {
		return err
	}
	// add new source metadata
	meta := &sandpiper.SliceMetadata{SliceID: sliceID}
	for _, m := range metaArray {
		meta.Key = m.Key
		meta.Value = m.Value
		if err := db.Insert(meta); err != nil {
			return err
		}
	}
	return nil
}

// Grains returns a list of grains for a slice (with brief or all fields)
// assumes allowed to do this
func (s *Sync) Grains(db orm.DB, sliceID uuid.UUID, briefFlag bool) ([]sandpiper.Grain, error) {
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
func (s *Sync) DeleteGrains(db orm.DB, ids []uuid.UUID) (err error) {
	if len(ids) > 0 {
		_, err = db.Model((*sandpiper.Grain)(nil)).Where("id in (?)", pg.In(ids)).Delete()
	}
	return err
}

// BeginSyncUpdate starts a quasi-transaction on a slice sync
func (s *Sync) BeginSyncUpdate(db orm.DB, sliceID uuid.UUID) error {
	m := sandpiper.Slice{
		ID:              sliceID,
		SyncStatus:      sandpiper.SyncStatusUpdating,
		LastSyncAttempt: time.Now(),
	}
	if _, err := db.Model(&m).WherePK().UpdateNotZero(); err != nil {
		return err
	}
	return nil
}

// FinalizeSyncUpdate completes the quasi-transaction on a slice sync
func (s *Sync) FinalizeSyncUpdate(db orm.DB, sliceID uuid.UUID, err error) error {
	var goodSync time.Time

	status := sandpiper.SyncStatusError
	if err == nil {
		status = sandpiper.SyncStatusSuccess
		goodSync = time.Now()
	}
	m := sandpiper.Slice{
		ID:           sliceID,
		SyncStatus:   status,
		LastGoodSync: goodSync,
	}
	if _, err := db.Model(&m).WherePK().UpdateNotZero(); err != nil {
		return err
	}
	return nil
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
