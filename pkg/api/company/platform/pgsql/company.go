// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package pgsql

// company service database access

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

// Company represents the client for company table
type Company struct{}

// NewCompany returns a new company model
func NewCompany() *Company {
	return &Company{}
}

// Custom errors
var (
	// ErrAlreadyExists indicates the company name is already used
	ErrAlreadyExists = echo.NewHTTPError(http.StatusInternalServerError, "Company name already exists.")
)

// Create creates a new company in database (assumes allowed to do this)
func (s *Company) Create(db orm.DB, company sandpiper.Company) (*sandpiper.Company, error) {
	// don't add if duplicate name
	if err := checkDuplicate(db, company.Name); err != nil {
		return nil, err
	}
	if err := db.Insert(&company); err != nil {
		return nil, err
	}
	return &company, nil
}

// View returns a single company by ID (assumes allowed to do this)
func (s *Company) View(db orm.DB, id uuid.UUID) (*sandpiper.Company, error) {
	var company = &sandpiper.Company{ID: id}

	err := db.Model(company).Relation("Users").WherePK().Select()
	if err != nil {
		return nil, err
	}
	return company, nil
}

// List returns list of all companies
func (s *Company) List(db orm.DB, sc *sandpiper.Scope, p *params.Params) (companies []sandpiper.Company, err error) {

	q := db.Model(&companies).Relation("Users").Limit(p.Paging.Limit).Offset(p.Paging.Offset()).Order("name")
	if sc != nil {
		q.Where(sc.Condition, sc.ID)
	}

	p.Paging.Count, err = q.SelectAndCount()
	if err != nil {
		return nil, err
	}

	return companies, nil
}

// Update updates company info by primary key (assumes allowed to do this)
func (s *Company) Update(db orm.DB, company *sandpiper.Company) error {
	_, err := db.Model(company).WherePK().UpdateNotZero()
	return err
}

// Delete removes the company (and any related subscriptions)
// All related users must be deleted elsewhere first
func (s *Company) Delete(db orm.DB, company *sandpiper.Company) error {
	return db.Delete(company)
}

// Server returns a single active company by ID for the sync process
func (s *Company) Server(db orm.DB, id uuid.UUID) (*sandpiper.Company, error) {
	company := &sandpiper.Company{ID: id}
	err := db.Model(company).Column("id", "name", "sync_addr").WherePK().Where("active = true").Select()
	if err != nil {
		return nil, err
	}
	return company, nil
}

// Servers returns a list of active companies (except ours) for the sync process optionally limited by "name"
func (s *Company) Servers(db orm.DB, ourCompanyID uuid.UUID, name string) ([]sandpiper.Company, error) {
	var companies []sandpiper.Company

	q := db.Model(&companies).
		Column("id", "name", "sync_addr", "sync_api_key", "active").
		Where("id <> ?", ourCompanyID).Where("sync_addr <> ''").Where("active = true")
	if name != "" {
		q = q.Where("lower(name) = ?", strings.ToLower(name))
	}
	err := q.Order("name").Select()
	if err != nil {
		return nil, err
	}
	return companies, nil
}

// checkDuplicate returns true if name found in database
func checkDuplicate(db orm.DB, name string) error {
	// attempt to select by unique key
	m := new(sandpiper.Company)
	err := db.Model(m).
		Column("id").
		Where("lower(name) = ?", strings.ToLower(name)).
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
