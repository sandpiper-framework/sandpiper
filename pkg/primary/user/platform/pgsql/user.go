// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql

// user service database access

import (
	"net/http"
	"strings"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/shared/model"
)

// User represents the client for user table
type User struct{}

// NewUser returns a new user database instance
func NewUser() *User {
	return &User{}
}

// Custom errors
var (
	// ErrAlreadyExists the username or email already exists
	ErrAlreadyExists = echo.NewHTTPError(http.StatusInternalServerError, "Username or email already exists.")
)

// Create creates a new user in the database (id is serially assigned)
func (u *User) Create(db orm.DB, usr sandpiper.User) (*sandpiper.User, error) {
	// make sure we can add this user
	if err := checkDuplicate(db, usr.Username, usr.Email); err != nil {
		return nil, err
	}

	if err := db.Insert(&usr); err != nil {
		return nil, err
	}
	return &usr, nil
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

// Update updates user's contact info
func (u *User) Update(db orm.DB, user *sandpiper.User) error {
	_, err := db.Model(user).UpdateNotZero()
	return err
}

// List returns list of all users retrievable for the current user, depending on role
func (u *User) List(db orm.DB, sc *sandpiper.Scope, p *sandpiper.Pagination) ([]sandpiper.User, error) {
	var users []sandpiper.User

	q := db.Model(&users).Limit(p.Limit).Offset(p.Offset).Order("username")
	if sc != nil {
		q.Where(sc.Condition, sc.ID)
	}
	if err := q.Select(); err != nil {
		return nil, err
	}
	return users, nil
}

// Delete permanently removes a user record from the database (not just making inactive)
func (u *User) Delete(db orm.DB, user *sandpiper.User) error {
	return db.Delete(user)
}

// checkDuplicate returns true if name found in database
func checkDuplicate(db orm.DB, name, email string) error {
	// attempt to select by unique keys
	m := new(sandpiper.User)
	err := db.Model(m).
		Column("id").
		Where("lower(username) = ? or lower(email) = ?", strings.ToLower(name), strings.ToLower(email)).
		Select()

	switch err {
	case pg.ErrNoRows: // ok to add
		return nil
	case nil: // found a row, so a duplicate
		return ErrAlreadyExists
	default: // return any other problem found
		return err
	}
}