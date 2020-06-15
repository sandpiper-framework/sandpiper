// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package company contains services for the companies resource. Companies have
// users and subscriptions to slices. Users must be a "company admin" to make changes
// to their company information.
package company

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/params"

	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
)

// Create adds a new company if administrator
func (s *Company) Create(c echo.Context, req sandpiper.Company) (*sandpiper.Company, error) {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return nil, err
	}
	return s.sdb.Create(s.db, req)
}

// List returns list of companies that you can view
func (s *Company) List(c echo.Context, p *params.Params) ([]sandpiper.Company, error) {
	q, err := s.rbac.EnforceScope(c)
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
	ID         uuid.UUID
	Name       string
	SyncAddr   string
	SyncAPIKey string
	SyncUserID int
	Active     bool
}

// Update updates company information
func (s *Company) Update(c echo.Context, r *Update) (*sandpiper.Company, error) {
	if err := s.rbac.EnforceCompany(c, r.ID); err != nil {
		return nil, err
	}
	company := &sandpiper.Company{
		ID:         r.ID,
		Name:       r.Name,
		SyncAddr:   r.SyncAddr,
		SyncAPIKey: r.SyncAPIKey,
		SyncUserID: r.SyncUserID,
		Active:     r.Active,
	}
	err := s.sdb.Update(s.db, company)
	if err != nil {
		return nil, err
	}
	return s.sdb.View(s.db, r.ID)
}

// Server returns a single server (company) that you can sync (if  admin)
func (s *Company) Server(c echo.Context, companyID uuid.UUID) (*sandpiper.Company, error) {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return nil, err
	}
	return s.sdb.Server(s.db, companyID)
}

// Servers returns list of servers (companies) that you can sync (if  admin)
func (s *Company) Servers(c echo.Context, name string) ([]sandpiper.Company, error) {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return nil, err
	}
	u := s.rbac.CurrentUser(c)
	return s.sdb.Servers(s.db, u.CompanyID, name)
}
