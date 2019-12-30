// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package company contains services for the companies resource
package company

import (
	"github.com/labstack/echo/v4"
	"github.com/google/uuid"

	"autocare.org/sandpiper/internal/model"
	"autocare.org/sandpiper/internal/scope"
)

// Create creates a new company if allowed
func (s *Company) Create(c echo.Context, req sandpiper.Company) (*sandpiper.Company, error) {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return nil, err
	}
	return s.sdb.Create(s.db, req)
}

// List returns list of companies that you can view
func (s *Company) List(c echo.Context, p *sandpiper.Pagination) ([]sandpiper.Company, error) {
	au := s.rbac.CurrentUser(c)
	q, err := scope.Limit(au)
	if err != nil {
		return nil, err
	}
	return s.sdb.List(s.db, q, p)
}

// View returns a single company if allowed
func (s *Company) View(c echo.Context, companyID uuid.UUID) (*sandpiper.Company, error) {
	if err := s.rbac.EnforceCompany(c, companyID); err != nil {
		return nil, err
	}
	return s.sdb.View(s.db, companyID)
}

// Delete deletes a company if allowed
func (s *Company) Delete(c echo.Context, id uuid.UUID) error {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return err
	}
	company, err := s.sdb.View(s.db, id)
	if err != nil {
		return err
	}
	return s.sdb.Delete(s.db, company)
}

// Update contains company's information used for updating
type Update struct {
	ID     uuid.UUID
	Name   string
	Active bool
}

// Update updates company information
func (s *Company) Update(c echo.Context, r *Update) (*sandpiper.Company, error) {

	if err := s.rbac.EnforceCompany(c, r.ID); err != nil {
		return nil, err
	}

	company := &sandpiper.Company{
		ID:     r.ID,
		Name:   r.Name,
		Active: r.Active,
	}
	err := s.sdb.Update(s.db, company)
	if err != nil {
		return nil, err
	}

	return s.sdb.View(s.db, r.ID)
}
