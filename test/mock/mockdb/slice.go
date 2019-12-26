package mockdb

import (
	"github.com/go-pg/pg/v9/orm"

	"autocare.org/sandpiper/internal/model"
)

// Slice database mock
type Slice struct {
	CreateFn         func(orm.DB, sandpiper.User) (*sandpiper.User, error)
	ViewFn           func(orm.DB, int) (*sandpiper.User, error)
	FindByNameFn     func(orm.DB, string) (*sandpiper.User, error)
	ListFn           func(orm.DB, *sandpiper.ListQuery, *sandpiper.Pagination) ([]sandpiper.User, error)
	DeleteFn         func(orm.DB, *sandpiper.User) error
	UpdateFn         func(orm.DB, *sandpiper.User) error
}

// Create mock
func (s *Slice) Create(db orm.DB, usr sandpiper.User) (*sandpiper.User, error) {
	return s.CreateFn(db, usr)
}

// View mock
func (s *Slice) View(db orm.DB, id int) (*sandpiper.User, error) {
	return s.ViewFn(db, id)
}

// FindByName mock
func (s *Slice) FindByUsername(db orm.DB, uname string) (*sandpiper.User, error) {
	return s.FindByNameFn(db, uname)
}

// List mock
func (s *Slice) List(db orm.DB, lq *sandpiper.ListQuery, p *sandpiper.Pagination) ([]sandpiper.User, error) {
	return s.ListFn(db, lq, p)
}

// Delete mock
func (s *Slice) Delete(db orm.DB, usr *sandpiper.User) error {
	return s.DeleteFn(db, usr)
}

// Update mock
func (s *Slice) Update(db orm.DB, usr *sandpiper.User) error {
	return s.UpdateFn(db, usr)
}
