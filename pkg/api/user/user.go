// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package user contains user application services
package user

import (
	"github.com/labstack/echo/v4"

	"sandpiper/pkg/shared/model"
	"sandpiper/pkg/shared/secure"
)

// Create creates a new user account
func (u *User) Create(c echo.Context, req sandpiper.User) (*sandpiper.User, error) {
	if err := u.rbac.AccountCreate(c, req.Role, req.CompanyID); err != nil {
		return nil, err
	}
	req.Password = u.sec.Hash(req.Password)
	return u.sdb.Create(u.db, req)
}

// List returns list of users
func (u *User) List(c echo.Context, p *sandpiper.Pagination) ([]sandpiper.User, error) {
	q, err := u.rbac.EnforceScope(c)
	if err != nil {
		return nil, err
	}
	return u.sdb.List(u.db, q, p)
}

// View returns a single user if allowed
func (u *User) View(c echo.Context, id int) (*sandpiper.User, error) {
	if err := u.rbac.EnforceUser(c, id); err != nil {
		return nil, err
	}
	return u.sdb.View(u.db, id)
}

// Delete deletes a user
func (u *User) Delete(c echo.Context, id int) error {
	user, err := u.sdb.View(u.db, id)
	if err != nil {
		return err
	}
	if err := u.rbac.IsLowerRole(c, user.Role); err != nil {
		return err
	}
	return u.sdb.Delete(u.db, user)
}

// Update contains user's information used for updating
type Update struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Active    bool
}

// Update updates user's contact information
func (u *User) Update(c echo.Context, r *Update) (*sandpiper.User, error) {
	if err := u.rbac.EnforceUser(c, r.ID); err != nil {
		return nil, err
	}

	if err := u.sdb.Update(u.db, &sandpiper.User{
		ID:        r.ID,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
		Phone:     r.Phone,
	}); err != nil {
		return nil, err
	}

	return u.sdb.View(u.db, r.ID)
}

// CreateAPIKey creates a sync user for a company (if necessary) and generates an apikey
func (u *User) CreateAPIKey(c echo.Context) (*sandpiper.APIKey, error) {
	// must be an api call on primary server
	if err := u.rbac.EnforceServerRole(sandpiper.PrimaryServer); err != nil {
		return nil, err
	}
	// caller must be companyAdmin
	if err := u.rbac.EnforceRole(c, sandpiper.CompanyAdminRole); err != nil {
		return nil, err
	}
	// get the company's sync user (or create one)
	companyID := u.rbac.CurrentUser(c).CompanyID
	usr, err := u.sdb.CompanySyncUser(u.db, companyID)
	if err != nil {
		return nil, err
	}

	// generate a new plain-text password for the sync_user and update the user record
	pw, err := u.sec.RandomPassword(26)
	if err != nil {
		return nil, err
	}
	usr.ChangePassword(u.sec.Hash(pw))
	if err := u.sdb.UpdateSyncUser(u.db, usr); err != nil {
		return nil, err
	}

	// encrypt these credentials in an api_key
	creds := &secure.Credentials{
		Username: usr.Username,
		Password: pw,
	}
	key, err := creds.APIKey(u.sec.APIKeySecret())
	if err != nil {
		return nil, err
	}

	return &sandpiper.APIKey{PrimaryID: u.rbac.OurServer().ID, SyncAPIKey: string(key)}, err
}
