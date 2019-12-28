// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql

// company service database access

import (
	"net/http"
	"strings"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"
	"github.com/google/uuid"

	"autocare.org/sandpiper/internal/model"
)

// Company represents the client for company table
type Company struct{}

// NewCompany returns a new company database instance
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
	// don't add if the name already exists
	if nameExists(db, company.Name) {
		return nil, ErrAlreadyExists
	}
	if err := db.Insert(&company); err != nil {
		return nil, err
	}
	return &company, nil
}

// View returns a single company by ID (assumes allowed to do this)
func (s *Company) View(db orm.DB, id uuid.UUID) (*sandpiper.Company, error) {
	var company = &sandpiper.Company{ID: id}

	err := db.Select(company)
	if err != nil {
		return nil, err
	}
	return company, nil
}

// Update updates company info by primary key (assumes allowed to do this)
func (s *Company) Update(db orm.DB, company *sandpiper.Company) error {
	_, err := db.Model(company).Update()
	return err
}

// List returns list of all companies
func (s *Company) List(db orm.DB, qp *sandpiper.ListQuery, p *sandpiper.Pagination) ([]sandpiper.Company, error) {
	var companies []sandpiper.Company

	q := db.Model(&companies).Limit(p.Limit).Offset(p.Offset).Where("deleted_at is null").Order("name")
	if qp != nil {
		q.Where(qp.Query, qp.ID)
	}
	if err := q.Select(); err != nil {
		return nil, err
	}
	return companies, nil
}

// Delete sets deleted_at for a company
func (s *Company) Delete(db orm.DB, company *sandpiper.Company) error {
	return db.Delete(company)
}

// nameExists returns true if name found in database
func nameExists(db orm.DB, name string) bool {
	m := new(sandpiper.Company)
	err := db.Model(m).Where("lower(name) = ?", strings.ToLower(name)).Select()
	return err == pg.ErrNoRows
}
