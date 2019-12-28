// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package user contains user application services
package user

import (
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/internal/model"
	"autocare.org/sandpiper/internal/query"
)

// Create creates a new user account
func (u *User) Create(c echo.Context, req sandpiper.User) (*sandpiper.User, error) {
	if err := u.rbac.AccountCreate(c, req.Role, req.CompanyID); err != nil {
		return nil, err
	}
	req.Password = u.sec.Hash(req.Password)
	return u.udb.Create(u.db, req)
}

// List returns list of users
func (u *User) List(c echo.Context, p *sandpiper.Pagination) ([]sandpiper.User, error) {
	au := u.rbac.CurrentUser(c)
	q, err := query.List(au)
	if err != nil {
		return nil, err
	}
	return u.udb.List(u.db, q, p)
}

// View returns a single user if allowed
func (u *User) View(c echo.Context, id int) (*sandpiper.User, error) {
	if err := u.rbac.EnforceUser(c, id); err != nil {
		return nil, err
	}
	return u.udb.View(u.db, id)
}

// Delete deletes a user
func (u *User) Delete(c echo.Context, id int) error {
	user, err := u.udb.View(u.db, id)
	if err != nil {
		return err
	}
	if err := u.rbac.IsLowerRole(c, user.Role); err != nil {
		return err
	}
	return u.udb.Delete(u.db, user)
}

// Update contains user's information used for updating
type Update struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

// Update updates user's contact information
func (u *User) Update(c echo.Context, r *Update) (*sandpiper.User, error) {
	if err := u.rbac.EnforceUser(c, r.ID); err != nil {
		return nil, err
	}

	if err := u.udb.Update(u.db, &sandpiper.User{
		ID:        r.ID,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
		Phone:     r.Phone,
	}); err != nil {
		return nil, err
	}

	return u.udb.View(u.db, r.ID)
}
