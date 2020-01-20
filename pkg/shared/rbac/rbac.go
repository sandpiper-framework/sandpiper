// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package rbac implements a role based access control mechanism as a service.
package rbac

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/shared/model"
)

// New creates new RBAC service
func New() *Service {
	return &Service{
		ScopingField: "company_id", // default, but is just "id" for company model
	}
}

// Service is RBAC enforcement service
type Service struct {
	ScopingField string // company id field name for scoping
}

// CurrentUser returns login data stored in jwt token
func (s *Service) CurrentUser(c echo.Context) *sandpiper.AuthUser {
	return &sandpiper.AuthUser{
		ID:        c.Get("id").(int),
		Username:  c.Get("username").(string),
		CompanyID: c.Get("company_id").(uuid.UUID),
		Email:     c.Get("email").(string),
		Role:      c.Get("role").(sandpiper.AccessLevel),
	}
}

// AccountCreate performs auth check when creating a new account
func (s *Service) AccountCreate(c echo.Context, role sandpiper.AccessLevel, companyID uuid.UUID) error {
	if err := s.EnforceCompany(c, companyID); err != nil {
		return err
	}
	return s.IsLowerRole(c, role)
}

// EnforceRole authorizes request by AccessLevel
func (s *Service) EnforceRole(c echo.Context, r sandpiper.AccessLevel) error {
	return checkBool(!(c.Get("role").(sandpiper.AccessLevel) > r))
}

// EnforceUser checks whether the request to change user data is done by the same user
func (s *Service) EnforceUser(c echo.Context, ID int) error {
	// TODO: Implement querying db and checking the requested user's company_id (??)
	// to allow company admins to view the user
	if s.isAdmin(c) {
		return nil
	}
	return checkBool(c.Get("id").(int) == ID)
}

// EnforceCompany checks whether the request to apply change to company data
// is done by the user belonging to the that company and that the user has role CompanyAdmin.
// If user has admin role, the check for company doesnt need to pass.
func (s *Service) EnforceCompany(c echo.Context, ID uuid.UUID) error {
	if s.isAdmin(c) {
		return nil
	}
	if err := s.EnforceRole(c, sandpiper.CompanyAdminRole); err != nil {
		return err
	}
	return checkBool(c.Get("company_id").(uuid.UUID) == ID)
}

// EnforceSubscription makes sure we can do something with this (populated) subscription
func (s *Service) EnforceSubscription(c echo.Context, sub sandpiper.Subscription) error {
	// allow admin or company admin with matching company
	return s.EnforceCompany(c, sub.CompanyID)
}

// EnforceScope uses the current user to determine if scoping needs to be added to a query
func (s *Service) EnforceScope(c echo.Context) (*sandpiper.Scope, error) {
	au := s.CurrentUser(c)
	return au.ApplyScope(s.ScopingField)
}

// IsLowerRole checks whether the requesting user has higher role than the user it wants to change
// Used for account creation/deletion
func (s *Service) IsLowerRole(c echo.Context, r sandpiper.AccessLevel) error {
	return checkBool(c.Get("role").(sandpiper.AccessLevel) < r)
}

func (s *Service) isAdmin(c echo.Context) bool {
	return !(c.Get("role").(sandpiper.AccessLevel) > sandpiper.AdminRole)
}

func (s *Service) isCompanyAdmin(c echo.Context) bool {
	return !(c.Get("role").(sandpiper.AccessLevel) > sandpiper.CompanyAdminRole)
}

func checkBool(b bool) error {
	if b {
		return nil
	}
	return echo.ErrForbidden
}
