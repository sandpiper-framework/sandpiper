package pgsql

// slice service database access

import (
	"net/http"
	"strings"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/internal/model"
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

// Create creates a new slice in database
func (s *Slice) Create(db orm.DB, slice sandpiper.Slice) (*sandpiper.Slice, error) {
	var dummy = new(sandpiper.Slice)

	// don't add if the name already exists
	sliceName := strings.ToLower(slice.Name)
	err := db.Model(dummy).Where("lower(name) = ? and deleted_at is null", sliceName).Select()
	if err != nil && err != pg.ErrNoRows {
		return nil, ErrAlreadyExists
	}

	if err := db.Insert(&slice); err != nil {
		return nil, err
	}
	return &slice, nil
}

// View returns single user by ID
func (s *Slice) View(db orm.DB, id sandpiper.UUID) (*sandpiper.Slice, error) {
	var slice = &sandpiper.Slice{ID: id}

	err := db.Select(slice)
	if err != nil {
		return nil, err
	}

	return slice, nil
}

// Update updates slice info by primary key
func (s *Slice) Update(db orm.DB, user *sandpiper.Slice) error {
	_, err := db.Model(user).Update()
	return err
}

// List returns list of all slices
func (s *Slice) List(db orm.DB, qp *sandpiper.ListQuery, p *sandpiper.Pagination) ([]sandpiper.Slice, error) {
	var slices []sandpiper.Slice

	q := db.Model(&slices).Limit(p.Limit).Offset(p.Offset).Where("deleted_at is null").Order("slice_name")
	if qp != nil {
		q.Where(qp.Query, qp.ID)
	}
	if err := q.Select(); err != nil {
		return nil, err
	}
	return slices, nil
}

// Delete sets deleted_at for a slice
func (s *Slice) Delete(db orm.DB, slice *sandpiper.Slice) error {
	return db.Delete(slice)
}
