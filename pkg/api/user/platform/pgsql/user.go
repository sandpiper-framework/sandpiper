// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package pgsql

// user service database access

import (
	"net/http"
	"strings"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/params"
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

// Update updates user info
func (u *User) Update(db orm.DB, user *sandpiper.User) error {
	_, err := db.Model(user).WherePK().UpdateNotZero()
	return err
}

// List returns all users retrievable by the current user, depending on role
func (u *User) List(db orm.DB, p *params.Params, sc *sandpiper.Scope) (users []sandpiper.User, err error) {

	q := db.Model(&users).Limit(p.Paging.Limit).Offset(p.Paging.Offset()).Order("username")
	if sc != nil {
		q.Where(sc.Condition, sc.ID)
	}
	p.AddSort(q)
	p.AddFilter(q)

	p.Paging.Count, err = q.SelectAndCount()
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Delete permanently removes a user record from the database (not just making inactive)
func (u *User) Delete(db orm.DB, user *sandpiper.User) error {
	return db.Delete(user)
}

// CompanySyncUser returns a company's sync_user. If one is not assigned, it creates one and links
// it to the company. Always change the username. Only need to return User.ID, User.Username.
func (u *User) CompanySyncUser(db orm.DB, companyID uuid.UUID) (*sandpiper.User, error) {
	usr := &sandpiper.User{
		Username: "sync_" + uuid.New().String(), // randomly unique, changed each time called
	}

	// get sync_user_id (on primary) for the secondary company requesting the key
	company := &sandpiper.Company{ID: companyID}
	err := db.Model(company).Column("sync_user_id").WherePK().Select()
	if err != nil {
		return nil, err
	}

	// If company doesn't have a sync_user, create one and assign it to the company
	if company.SyncUserID == 0 {
		// add the sync user without password
		su := &sandpiper.User{
			Username:  usr.Username,
			FirstName: "sync",
			LastName:  "sync",
			Email:     "sync",
			Active:    true,
			Role:      sandpiper.SyncRole,
			CompanyID: companyID,
		}
		if err := db.Insert(su); err != nil {
			return nil, err
		}
		// update company with new sync user (serially assigned)
		company.SyncUserID = su.ID
		_, err = db.Model(company).Column("sync_user_id").WherePK().Update()
		if err != nil {
			return nil, err
		}
	}
	usr.ID = company.SyncUserID

	return usr, err
}

// UpdateSyncUser updates database with changed password data from the model values (by primary key)
func (u *User) UpdateSyncUser(db orm.DB, user *sandpiper.User) error {
	_, err := db.Model(user).Column("username", "password", "password_changed").WherePK().Update()
	return err
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
