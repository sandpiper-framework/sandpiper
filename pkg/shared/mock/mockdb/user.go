package mockdb

import (
	"github.com/go-pg/pg/v9/orm"

	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
)

// User database mock
type User struct {
	CreateFn         func(orm.DB, sandpiper.User) (*sandpiper.User, error)
	ViewFn           func(orm.DB, int) (*sandpiper.User, error)
	FindByUsernameFn func(orm.DB, string) (*sandpiper.User, error)
	FindByTokenFn    func(orm.DB, string) (*sandpiper.User, error)
	ListFn           func(orm.DB, *sandpiper.Scope, *sandpiper.Pagination) ([]sandpiper.User, error)
	DeleteFn         func(orm.DB, *sandpiper.User) error
	UpdateFn         func(orm.DB, *sandpiper.User) error
}

// Create mock
func (u *User) Create(db orm.DB, usr sandpiper.User) (*sandpiper.User, error) {
	return u.CreateFn(db, usr)
}

// View mock
func (u *User) View(db orm.DB, id int) (*sandpiper.User, error) {
	return u.ViewFn(db, id)
}

// FindByUsername mock
func (u *User) FindByUsername(db orm.DB, uname string) (*sandpiper.User, error) {
	return u.FindByUsernameFn(db, uname)
}

// FindByToken mock
func (u *User) FindByToken(db orm.DB, token string) (*sandpiper.User, error) {
	return u.FindByTokenFn(db, token)
}

// List mock
func (u *User) List(db orm.DB, lq *sandpiper.Scope, p *sandpiper.Pagination) ([]sandpiper.User, error) {
	return u.ListFn(db, lq, p)
}

// Delete mock
func (u *User) Delete(db orm.DB, usr *sandpiper.User) error {
	return u.DeleteFn(db, usr)
}

// Update mock
func (u *User) Update(db orm.DB, usr *sandpiper.User) error {
	return u.UpdateFn(db, usr)
}
