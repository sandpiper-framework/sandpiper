// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql

// sync service database access

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"

	"sandpiper/pkg/shared/model"
)

var (
	// ErrAlreadyExists indicates that the name already exists
	ErrAlreadyExists = echo.NewHTTPError(http.StatusInternalServerError, "Subscription name already exists.")
)

// Sync represents the client for sync table
type Sync struct{}

// NewSync returns a new sync database instance
func NewSync() *Sync {
	return &Sync{}
}

// LogActivity adds a sync log entry to the activity table
func (s *Sync) LogActivity(db orm.DB, subID uuid.UUID, slice *sandpiper.Slice, d time.Duration) error {
	msg := "synced " + slice.Name
	if !slice.AllowSync {
		msg = "locked on primary"
	}
	// log this event in our activity table
	activity := sandpiper.Activity{
		SubID:   subID,
		Success: slice.AllowSync,
		Message: msg,
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

// Subscriptions returns list of all local subscriptions for a company
func (s *Sync) Subscriptions(db orm.DB, companyID uuid.UUID) ([]sandpiper.Subscription, error) {
	var subs []sandpiper.Subscription

	err := db.Model(&subs).Where("company_id = ?", companyID).Select()
	if err != nil {
		return nil, err
	}
	return subs, nil
}

// AddSubscription creates a new Subscription in database
func (s *Sync) AddSubscription(db orm.DB, sub sandpiper.Subscription) error {
	// make sure the name is not a duplicate
	if err := checkDuplicate(db, sub.Name); err != nil {
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

// checkDuplicate returns true if name found in database
func checkDuplicate(db orm.DB, name string) error {
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
