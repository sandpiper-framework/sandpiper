// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql

// grain service database access

import (
	"net/http"

	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/internal/model"
	"autocare.org/sandpiper/internal/scope"
)

// Grain represents the client for grain table
type Grain struct{}

// NewGrain returns a new grain database instance
func NewGrain() *Grain {
	return &Grain{}
}

// Custom errors
var (
	// ErrAlreadyExists indicates the grain name is already used
	ErrAlreadyExists = echo.NewHTTPError(http.StatusInternalServerError, "Grain name already exists.")
)

// Create creates a new grain in database (assumes allowed to do this)
func (s *Grain) Create(db orm.DB, grain sandpiper.Grain) (*sandpiper.Grain, error) {

	if err := db.Insert(&grain); err != nil {
		return nil, err
	}
	return &grain, nil
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
		Relation("slices._").
		Relation("subscriptions._").
		Where("slice_id = ? and subscriber_id = ?", sliceID, companyID).
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
		Relation("slices._").
		Limit(p.Limit).Offset(p.Offset)
	if sc != nil {
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
