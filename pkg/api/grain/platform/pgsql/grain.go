// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package pgsql

// grain service database access

// The `grain` contains the actual data `payload` that is exchanged between trading partners.
// The payload is transferred (and stored) as a (possibly encoded) binary object.

import (
	"github.com/sandpiper-framework/sandpiper/pkg/shared/params"
	"net/http"
	"strings"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
)

// Grain represents the client for grain table
type Grain struct{}

// NewGrain returns a new grain database instance
func NewGrain() *Grain {
	return &Grain{}
}

// Custom errors
var (
	// ErrGrainNotFound indicates select returned no rows
	ErrGrainNotFound = echo.NewHTTPError(http.StatusNotFound, "Grain does not exist.")
)

// Create creates a new grain in database (assumes allowed to do this).
func (s *Grain) Create(db orm.DB, replaceFlag bool, grain *sandpiper.Grain) (*sandpiper.Grain, error) {
	// key is always lowercase to allow faster lookups without a function index
	grain.Key = strings.ToLower(grain.Key)

	if replaceFlag {
		if err := removeExistingGrain(db, *grain.SliceID, grain.Key); err != nil {
			return nil, err
		}
	}

	if err := db.Insert(grain); err != nil {
		return nil, err
	}
	return grain, nil
}

// View returns a single grain by ID (assumes allowed to do this)
func (s *Grain) View(db orm.DB, id uuid.UUID) (*sandpiper.Grain, error) {
	var grain = &sandpiper.Grain{ID: id}

	err := db.Model(grain).
		Column("grain.id", "slice_id", "grain_key", "source", "encoding", "payload", "grain.created_at").
		ColumnExpr("length(payload) AS payload_len").
		Relation("Slice").WherePK().Select()
	if err != nil {
		return nil, selectError(err)
	}
	return grain, nil
}

// ViewByKeys returns minimal grain information if found, an empty grain if not found
func (s *Grain) ViewByKeys(db orm.DB, sliceID uuid.UUID, grainKey string, payloadFlag bool) (*sandpiper.Grain, error) {
	// columns to select (optionally returning payload)
	cols := "id, slice_id, grain_key, source, encoding, created_at, length(payload) AS payload_len"
	if payloadFlag {
		cols = cols + ", payload"
	}

	grain := new(sandpiper.Grain)
	err := db.Model(grain).ColumnExpr(cols).
		Where("slice_id = ? and grain_key = ?", sliceID, grainKey).
		Select()
	if err != nil && err != pg.ErrNoRows {
		return nil, err
	}
	return grain, nil
}

// CompanySubscribed checks if grain is included in a company's subscriptions.
func (s *Grain) CompanySubscribed(db orm.DB, companyID uuid.UUID, grainID uuid.UUID) bool {
	grain := new(sandpiper.Grain)
	err := db.Model(grain).Column("grain.id").
		Join("INNER JOIN subscriptions AS sub ON grain.slice_id = sub.slice_id").
		Where("sub.company_id = ?", companyID).
		Where("grain.id = ?", grainID).Select()
	if err == nil {
		return true
	}
	return false
}

// List returns a list of all grains with scoping and pagination (optionally for a slice)
func (s *Grain) List(db orm.DB, sliceID uuid.UUID, payloadFlag bool, sc *sandpiper.Scope, p *params.Params) ([]sandpiper.Grain, error) {
	var grains []sandpiper.Grain
	var q *orm.Query

	// columns to select (optionally returning payload)
	cols := "grain.id, grain.slice_id, grain_key, source, encoding, grain.created_at, length(payload) AS payload_len"
	if payloadFlag {
		cols = cols + ", payload"
	}

	// build the query
	switch {
	case sc != nil && sliceID != uuid.Nil:
		// both provided, join to subscriptions (by-passing slices table)
		q = db.Model(&grains).ColumnExpr(cols).
			Join("INNER JOIN subscriptions AS sub ON grain.slice_id = sub.slice_id").
			Where("sub.company_id = ?", sc.ID).Where("active = true")
	case sc != nil && sliceID == uuid.Nil:
		// provided scope without a slice
		// Use CTE query to get all "active" subscriptions for the scope (i.e. the company)
		q = db.Model((*sandpiper.Subscription)(nil)).
			Column("subscription.slice_id").Where(sc.Condition, sc.ID).Where("active = true").
			WrapWith("scope").Table("scope").
			Join("INNER JOIN grains AS grain ON grain.slice_id = scope.slice_id").
			ColumnExpr(cols)
	case sc == nil && sliceID != uuid.Nil:
		// provided slice without a scope, use simple where clause
		q = db.Model(&grains).ColumnExpr(cols).Where("slice_id = ?", sliceID)
	default:
		// neither provided, simply return all grains
		q = db.Model(&grains).ColumnExpr(cols)
	}

	// execute the query
	err := q.Limit(p.Paging.Limit).Offset(p.Paging.Offset()).Select(&grains)
	if err != nil {
		return nil, err
	}
	return grains, nil
}

// Delete permanently removes a grain by primary key (id)
func (s *Grain) Delete(db orm.DB, id uuid.UUID) error {
	grain := sandpiper.Grain{ID: id}
	return db.Delete(&grain)
}

// removeExistingGrain will remove a grain by alternate unique key. Only return real errors.
func removeExistingGrain(db orm.DB, sliceID uuid.UUID, grainKey string) error {
	// attempt to delete by unique keys
	m := new(sandpiper.Grain)
	_, err := db.Model(m).Where("slice_id = ? AND grain_key = ?", sliceID, grainKey).Delete()
	if err != nil && err != pg.ErrNoRows {
		return err
	}
	return nil
}

func selectError(err error) error {
	if err == pg.ErrNoRows {
		return ErrGrainNotFound
	}
	return err
}
