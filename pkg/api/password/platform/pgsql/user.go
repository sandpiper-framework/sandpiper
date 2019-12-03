package pgsql

// password service database access

import (
	"github.com/go-pg/pg/v9/orm"

	"autocare.org/sandpiper/pkg/model"
)

// User represents the client for user table
type User struct{}

// NewUser returns a new user database instance
func NewUser() *User {
	return &User{}
}

// View returns single user by ID
func (u *User) View(db orm.DB, id int) (*sandpiper.User, error) {
	user := &sandpiper.User{Base: sandpiper.Base{ID: id}}
	err := db.Select(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Update updates user's info
func (u *User) Update(db orm.DB, user *sandpiper.User) error {
	return db.Update(user)
}
