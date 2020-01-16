// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package company contains services for the companies resource. Companies have
// users and subscriptions to slices. Users must be a "company admin" to make changes
// to their company information.
package company

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/internal/model"
	"autocare.org/sandpiper/pkg/internal/scope"
)

// Create adds a new company if administrator
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

// Delete deletes a company if administrator
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

// Update contains company field request used for updating
type Update struct {
	ID       uuid.UUID
	Name     string
	SyncAddr string
	Active   bool
}

// Update updates company information
func (s *Company) Update(c echo.Context, r *Update) (*sandpiper.Company, error) {
	if err := s.rbac.EnforceCompany(c, r.ID); err != nil {
		return nil, err
	}
	company := &sandpiper.Company{
		ID:       r.ID,
		Name:     r.Name,
		SyncAddr: r.SyncAddr,
		Active:   r.Active,
	}
	err := s.sdb.Update(s.db, company)
	if err != nil {
		return nil, err
	}
	return s.sdb.View(s.db, r.ID)
}
