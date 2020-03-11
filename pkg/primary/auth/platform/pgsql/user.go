// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql

// auth service database access

import (
	"github.com/go-pg/pg/v9/orm"

	"autocare.org/sandpiper/pkg/shared/model"
)

// User represents the client for user table
type User struct{}

// NewUser returns a new user database instance (for the service)
func NewUser() *User {
	return &User{}
}

// View returns single user by ID
func (u *User) View(db orm.DB, id int) (*sandpiper.User, error) {
	var user = &sandpiper.User{ID: id}

	err := db.Select(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindByUsername queries for single user by username
func (u *User) FindByUsername(db orm.DB, uname string) (*sandpiper.User, error) {
	var user = new(sandpiper.User)

	err := db.Model(user).Where("username = ?", uname).Select()
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindByToken queries for single user by token
func (u *User) FindByToken(db orm.DB, token string) (*sandpiper.User, error) {
	var user = new(sandpiper.User)

	err := db.Model(user).Where("token = ?", token).Select()
	if err != nil {
		return nil, err
	}
	return user, err
}

// Update updates user's info
func (u *User) Update(db orm.DB, user *sandpiper.User) error {
	return db.Update(user)
}
