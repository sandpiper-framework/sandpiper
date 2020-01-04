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
	// don't add if the name already exists
	if nameExists(db, slice.Name) {
		return nil, ErrAlreadyExists
	}
	if err := db.Insert(&slice); err != nil {
		return nil, err
	}

	// insert any meta data
	meta := sandpiper.SliceMetaData{SliceID: slice.ID}
	for k, v := range slice.MetaData {
		meta.Key = k
		meta.Value = v
		if err := db.Insert(meta); err != nil {
			return nil, err
		}
	}

	return &slice, nil
}

// View returns a single slice with metadata by ID (assumes allowed to do this)
func (s *Slice) View(db orm.DB, id uuid.UUID) (*sandpiper.Slice, error) {
	var slice = &sandpiper.Slice{ID: id}

	/*	err := db.Model(company).
		ColumnExpr("company.*").
		ColumnExpr("users.*").
		Join("LEFT JOIN users").
		JoinOn("company.id = users.company_id").
		WherePK().First()
	*/

	err := db.Select(slice)
	if err != nil {
		return nil, err
	}

	return slice, nil
}

// ViewBySub returns a single slice by ID if included in company subscriptions.
func (s *Slice) ViewBySub(db orm.DB, companyID uuid.UUID, sliceID uuid.UUID) (*sandpiper.Slice, error) {
	var slice = &sandpiper.Slice{ID: sliceID}

	err := db.Model(slice).
		Relation("subscriptions._").
		Where("slice_id = ? and subscriber_id = ?", sliceID, companyID).
		Select()

	if err != nil {
		return nil, err
	}

	return slice, nil
}

// List returns list of all slices
func (s *Slice) List(db orm.DB, sc *scope.Clause, p *sandpiper.Pagination) ([]sandpiper.Slice, error) {
	var slices []sandpiper.Slice

	q := db.Model(&slices).Limit(p.Limit).Offset(p.Offset).Order("slice_name")
	if sc != nil {
		q.Where(sc.Condition, sc.ID)
	}
	if err := q.Select(); err != nil {
		return nil, err
	}
	return slices, nil
}

// Update updates slice info by primary key (assumes allowed to do this)
func (s *Slice) Update(db orm.DB, slice *sandpiper.Slice) error {
	_, err := db.Model(slice).Update()
	return err
}

// Delete a slice
func (s *Slice) Delete(db orm.DB, slice *sandpiper.Slice) error {
	return db.Delete(slice)
}

// nameExists returns true if name found in database
func nameExists(db orm.DB, name string) bool {
	m := new(sandpiper.Slice)
	err := db.Model(m).
		Column("id","name").
		Where("lower(name) = ?", strings.ToLower(name)).
		Select()
	return err != pg.ErrNoRows
}
