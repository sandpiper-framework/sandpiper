package mockdb

import (
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"

	"autocare.org/sandpiper/pkg/shared/model"
)

// Slice database mock
type Slice struct {
	CreateFn    func(orm.DB, sandpiper.Slice) (*sandpiper.Slice, error)
	ViewFn      func(orm.DB, uuid.UUID) (*sandpiper.Slice, error)
	ViewBySubFn func(db orm.DB, companyID uuid.UUID, sliceID uuid.UUID) (*sandpiper.Slice, error)
	ListFn      func(orm.DB, *sandpiper.Scope, *sandpiper.Pagination) ([]sandpiper.Slice, error)
	DeleteFn    func(orm.DB, *sandpiper.Slice) error
	UpdateFn    func(orm.DB, *sandpiper.Slice) error
}

// Create mock
func (s *Slice) Create(db orm.DB, slice sandpiper.Slice) (*sandpiper.Slice, error) {
	return s.CreateFn(db, slice)
}

// View mock
func (s *Slice) View(db orm.DB, id uuid.UUID) (*sandpiper.Slice, error) {
	return s.ViewFn(db, id)
}

// ViewBySub mock
func (s *Slice) ViewBySub(db orm.DB, companyID uuid.UUID, sliceID uuid.UUID) (*sandpiper.Slice, error) {
	return s.ViewBySubFn(db, companyID, sliceID)
}

// List mock
func (s *Slice) List(db orm.DB, lq *sandpiper.Scope, p *sandpiper.Pagination) ([]sandpiper.Slice, error) {
	return s.ListFn(db, lq, p)
}

// Delete mock
func (s *Slice) Delete(db orm.DB, slice *sandpiper.Slice) error {
	return s.DeleteFn(db, slice)
}

// Update mock
func (s *Slice) Update(db orm.DB, slice *sandpiper.Slice) error {
	return s.UpdateFn(db, slice)
}
