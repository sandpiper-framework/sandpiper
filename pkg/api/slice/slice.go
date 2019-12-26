// Package slice contains services for slices
package slice

import (
	"github.com/labstack/echo/v4"
	"github.com/satori/go.uuid"
	"time"

	"autocare.org/sandpiper/internal/model"
	"autocare.org/sandpiper/internal/query"
)

// Create creates a new slice to hold data-objects
func (s *Slice) Create(c echo.Context, req sandpiper.Slice) (*sandpiper.Slice, error) {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return nil, err
	}
	return s.sdb.Create(s.db, req)
}

// List returns list of slices
func (s *Slice) List(c echo.Context, p *sandpiper.Pagination) ([]sandpiper.Slice, error) {
	au := s.rbac.CurrentUser(c)
	q, err := query.List(au)
	if err != nil {
		return nil, err
	}
	return s.sdb.List(s.db, q, p)
}

// View returns a single slice if allowed
func (s *Slice) View(c echo.Context, id uuid.UUID) (*sandpiper.Slice, error) {
	//if err := s.rbac.EnforceVisibility(c); err != nil {
	//	return nil, err
	//}
	return s.sdb.View(s.db, id)
}

// Delete deletes a slice if allowed
func (s *Slice) Delete(c echo.Context, id uuid.UUID) error {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return err
	}

	slice, err := s.sdb.View(s.db, id)
	if err != nil {
		return err
	}

	return s.sdb.Delete(s.db, slice)
}

// Update contains slice's information used for updating
type Update struct {
	ID           uuid.UUID
	Name         string
	ContentHash  string
	ContentCount uint
	LastUpdate   time.Time
}

// Update updates slice information
func (s *Slice) Update(c echo.Context, r *Update) (*sandpiper.Slice, error) {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return nil, err
	}

	slice := &sandpiper.Slice{
		ID:           r.ID,
		Name:         r.Name,
		ContentHash:  r.ContentHash,
		ContentCount: r.ContentCount,
		LastUpdate:   r.LastUpdate,
	}

	err := s.sdb.Update(s.db, slice)
	if err != nil {
		return nil, err
	}

	return s.sdb.View(s.db, r.ID)
}
