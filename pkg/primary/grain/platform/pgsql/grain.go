// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql

// grain service database access

// The `grain` contains the actual data `payload` that is exchanged between trading partners.
// The payload is transferred (and stored) as a (possibly encoded) binary object.

import (
	"net/http"
	"strings"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/shared/model"
)

// Grain represents the client for grain table
type Grain struct{}

// NewGrain returns a new grain database instance
func NewGrain() *Grain {
	return &Grain{}
}

// Custom errors
var (
	// ErrAlreadyExists indicates the grain-type-key constraint would fail
	ErrAlreadyExists = echo.NewHTTPError(http.StatusInternalServerError, "Grain type/key already exists for this Slice.")
)

// Create creates a new grain in database (assumes allowed to do this)
func (s *Grain) Create(db orm.DB, replaceFlag bool, grain sandpiper.Grain) (*sandpiper.Grain, error) {
	// key is always lowercase to allow faster lookups
	grain.Key = strings.ToLower(grain.Key)

	if replaceFlag {
		err := removeExistingGrain(db, *grain.SliceID, grain.Type, grain.Key)
		if err != nil {
			return nil, err
		}
	} else {
		// todo: do we need to check? maybe we can check the insert error for existing keys (?)
		err := assertAddGrain(db, *grain.SliceID, grain.Type, grain.Key)
		if err != nil {
			return nil, err
		}
	}

	if err := db.Insert(&grain); err != nil {
		return nil, err
	}
	return &grain, nil
}

// View returns a single grain by ID (assumes allowed to do this)
func (s *Grain) View(db orm.DB, id uuid.UUID) (*sandpiper.Grain, error) {
	var grain = &sandpiper.Grain{ID: id}

	err := db.Model(grain).
		Column("grain.id", "grain_type", "grain_key", "encoding", "payload", "grain.created_at").
		Relation("Slice").WherePK().Select()
	if err != nil {
		return nil, err
	}
	return grain, nil
}

// Exists returns minimal grain information if found
func (s *Grain) Exists(db orm.DB, sliceID uuid.UUID, grainType, grainKey string) (*sandpiper.Grain, error) {
	grain := new(sandpiper.Grain)
	err := db.Model(grain).Column("grain.id", "source").
		Where("slice_id = ? and grain_type = ? and grain_key = ?", sliceID, grainType, grainKey).
		Select()
	if err != nil {
		return nil, err
	}
	return grain, nil
}

// CompanySubscribed checks if grain is included in a company's subscriptions.
func (s *Grain) CompanySubscribed(db orm.DB, companyID uuid.UUID, grainID uuid.UUID) bool {
	grain := new(sandpiper.Grain)
	err := db.Model(grain).Column("grain.id").
		Join("JOIN slices AS sl ON grain.slice_id = sl.id").
		Join("JOIN subscriptions AS sub ON sl.id = sub.slice_id").
		Where("sub.company_id = ?", companyID).
		Where("grain.id = ?", grainID).Select()
	if err == nil {
		return true
	}
	return false
}

// List returns a list of all grains with scoping and pagination
func (s *Grain) List(db orm.DB, payload bool, sc *sandpiper.Scope, p *sandpiper.Pagination) ([]sandpiper.Grain, error) {
	var grains []sandpiper.Grain
	var err error

	// columns to select (optionally returning payload)
	cols := "grain.id, grain_type, grain_key, encoding, grain.created_at"
	if payload {
		cols = cols + ", payload"
	}

	if sc != nil {
		// Use CTE query to get all subscriptions for the scope (i.e. the company)
		err = db.Model((*sandpiper.Subscription)(nil)).
			Column("subscription.slice_id").
			Where(sc.Condition, sc.ID).
			WrapWith("scope").Table("scope").
			Join("JOIN grains AS grain ON grain.slice_id = scope.slice_id").
			ColumnExpr(cols).
			Limit(p.Limit).Offset(p.Offset).Select(&grains)
	} else {
		// simple case with no scoping
		err = db.Model(&grains).ColumnExpr(cols).Limit(p.Limit).Offset(p.Offset).Select()
	}
	if err != nil {
		return nil, err
	}
	return grains, nil
}

// Delete permanently removes a grain by primary key (id)
func (s *Grain) Delete(db orm.DB, id uuid.UUID) error {
	grain := sandpiper.Grain{ID: id}
	return db.Delete(grain)
}

// assertAddGrain makes sure we can add this grain
func assertAddGrain(db orm.DB, sliceID uuid.UUID, grainType string, grainKey string) error {
	// attempt to select by unique keys
	m := new(sandpiper.Grain)
	err := db.Model(m).
		Column("id", "slice_id", "grain_type", "grain_key").
		Where("slice_id = ? and grain_type = ? and grain_key = ?", sliceID, grainType, grainKey).
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

// removeExistingGrain will remove a grain by alternate unique key. Only return real errors.
func removeExistingGrain(db orm.DB, sliceID uuid.UUID, grainType string, grainKey string) error {
	// attempt to delete by unique keys
	m := new(sandpiper.Grain)
	_, err := db.Model(m).
		Where("slice_id = ? and grain_type = ? and grain_key = ?", sliceID, grainType, grainKey).
		Delete()

	// todo: fix this test. See what is actually returned when delete vs nothing to delete.
	switch err {
	case pg.ErrNoRows: // ok to add
		return nil
	case nil: // found a row, so a duplicate
		return ErrAlreadyExists
	default: // return any other problem found
		return err
	}
}