// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql

// grain service database access

import (
	"github.com/go-pg/pg/v9"
	"net/http"
	"strings"

	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/internal/model"
	"autocare.org/sandpiper/pkg/internal/scope"
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
func (s *Grain) Create(db orm.DB, grain *sandpiper.Grain) (*sandpiper.Grain, error) {

	// key is always lowercase to allow faster lookups
	grain.Key = strings.ToLower(grain.Key)

	if duplicate(db, grain.SliceID, grain.Type, grain.Key) {
		return nil, ErrAlreadyExists
	}

	if err := db.Insert(&grain); err != nil {
		return nil, err
	}
	return grain, nil
}

// View returns a single grain by ID (assumes allowed to do this)
func (s *Grain) View(db orm.DB, id uuid.UUID) (*sandpiper.Grain, error) {
	var grain = &sandpiper.Grain{ID: id}

	if err := db.Select(grain); err != nil {
		return nil, err
	}

	return grain, nil
}

// ViewBySlice returns a single grain by ID if included in the supplied slice.
func (s *Grain) ViewBySlice(db orm.DB, grainID uuid.UUID) (*sandpiper.Grain, error) {
	// todo: implement this
	panic("implement me")
}

// ViewBySub returns a single grain by ID if included in company subscriptions.
func (s *Grain) ViewBySub(db orm.DB, companyID uuid.UUID, sliceID uuid.UUID) (*sandpiper.Grain, error) {
	var grain = &sandpiper.Grain{ID: sliceID}

	err := db.Model(grain).
		Relation("slices").
		Relation("subscriptions.company_id").
		Where("slice_id = ? and company_id = ?", sliceID, companyID).
		Select()

	if err != nil {
		return nil, err
	}
	return grain, nil
}

// List returns a list of all grains with scoping and pagination
func (s *Grain) List(db orm.DB, sc *scope.Clause, p *sandpiper.Pagination) ([]sandpiper.Grain, error) {
	var grains []sandpiper.Grain

	q := db.Model(&grains).
		Relation("slices").
		Limit(p.Limit).Offset(p.Offset)
	if sc != nil {
		q.Relation("subscriptions.company_id")
		q.Where(sc.Condition, sc.ID)
	}
	if err := q.Select(); err != nil {
		return nil, err
	}
	return grains, nil
}

// Delete permanently removes a grain
func (s *Grain) Delete(db orm.DB, grain *sandpiper.Grain) error {
	return db.Delete(grain)
}

// duplicate returns true if grain type/key found in database for a slice
func duplicate(db orm.DB, sliceID uuid.UUID, grainType sandpiper.GrainType, grainKey string) bool {
	m := new(sandpiper.Grain)
	err := db.Model(m).Where("slice_id = ? and grain_type = ? and grain_key = ?", sliceID, grainType, strings.ToLower(grainKey)).
		Select()
	return err != pg.ErrNoRows
}
