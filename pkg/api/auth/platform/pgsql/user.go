// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package pgsql

// auth service database access

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"
	"net/http"

	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
)

var (
	// ErrAuthNotFound message is a bit obscure on purpose
	ErrAuthNotFound = echo.NewHTTPError(http.StatusNotFound, "Authentication Error.")
)

// User represents the client for user table
type User struct{}

// NewUser returns a new user database instance (for the service)
func NewUser() *User {
	return &User{}
}

// View returns single user by ID (used by the /me route)
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
		return nil, selectError(err)
	}
	return user, nil
}

// FindByToken queries for single user by token
func (u *User) FindByToken(db orm.DB, token string) (*sandpiper.User, error) {
	var user = new(sandpiper.User)

	err := db.Model(user).Where("token = ?", token).Select()
	if err != nil {
		return nil, selectError(err)
	}
	return user, err
}

// Update updates user's info
func (u *User) Update(db orm.DB, user *sandpiper.User) error {
	return db.Update(user)
}

func selectError(err error) error {
	if err == pg.ErrNoRows {
		return ErrAuthNotFound
	}
	return err
}
