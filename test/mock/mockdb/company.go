package mockdb

import (
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"

	"autocare.org/sandpiper/internal/model"
)

// Company database mock
type Company struct {
	CreateFn         func(orm.DB, sandpiper.Company) (*sandpiper.Company, error)
	ListFn           func(orm.DB, *sandpiper.ListQuery, *sandpiper.Pagination) ([]sandpiper.Company, error)
	ViewFn           func(orm.DB, uuid.UUID) (*sandpiper.Company, error)
	DeleteFn         func(orm.DB, *sandpiper.Company) error
	UpdateFn         func(orm.DB, *sandpiper.Company) error
}

// Create mock
func (s *Company) Create(db orm.DB, cpy sandpiper.Company) (*sandpiper.Company, error) {
	return s.CreateFn(db, cpy)
}

// List mock
func (s *Company) List(db orm.DB, lq *sandpiper.ListQuery, p *sandpiper.Pagination) ([]sandpiper.Company, error) {
	return s.ListFn(db, lq, p)
}

// View mock
func (s *Company) View(db orm.DB, id uuid.UUID) (*sandpiper.Company, error) {
	return s.ViewFn(db, id)
}

// Delete mock
func (s *Company) Delete(db orm.DB, cpy *sandpiper.Company) error {
	return s.DeleteFn(db, cpy)
}

// Update mock
func (s *Company) Update(db orm.DB, cpy *sandpiper.Company) error {
	return s.UpdateFn(db, cpy)
}
