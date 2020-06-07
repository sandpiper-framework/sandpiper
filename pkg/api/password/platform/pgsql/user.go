// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package pgsql

// password service database access

import (
	"github.com/go-pg/pg/v9/orm"

	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
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
	_, err := db.Model(user).Column("password", "password_changed").WherePK().Update()
	return err
}
