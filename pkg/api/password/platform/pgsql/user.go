// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql

// password service database access

import (
	"github.com/go-pg/pg/v9/orm"

	"autocare.org/sandpiper/internal/model"
)

// User represents the client for user table
type User struct{}

// NewUser returns a new user database instance
func NewUser() *User {
	return &User{}
}

// View returns single user by ID
func (u *User) View(db orm.DB, id int) (*sandpiper.User, error) {
	user := &sandpiper.User{ID: id}
	err := db.Select(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Update updates database with changed password data from the model values (by primary key)
func (u *User) Update(db orm.DB, user *sandpiper.User) error {
	//return db.Update(user)
	_, err := db.Model(user).Column("password", "last_password_change").WherePK().Update()
	return err
}
