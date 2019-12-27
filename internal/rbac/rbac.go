// Package rbac implements a role based access control mechanism as a service.
package rbac

import (
	"github.com/labstack/echo/v4"
	"github.com/satori/go.uuid"

	"autocare.org/sandpiper/internal/model"
)

// New creates new RBAC service
func New() *Service {
	return &Service{}
}

// Service is RBAC enforcement service
type Service struct{}

// CurrentUser returns login data stored in jwt token
func (s *Service) CurrentUser(c echo.Context) *sandpiper.AuthUser {
	return &sandpiper.AuthUser{
		ID:         c.Get("id").(int),
		Username:   c.Get("username").(string),
		CompanyID:  c.Get("company_id").(uuid.UUID),
		Email:      c.Get("email").(string),
		Role:       c.Get("role").(sandpiper.AccessRole),
	}
}

// EnforceRole authorizes request by AccessRole
func (s *Service) EnforceRole(c echo.Context, r sandpiper.AccessRole) error {
	return checkBool(!(c.Get("role").(sandpiper.AccessRole) > r))
}

// EnforceUser checks whether the request to change user data is done by the same user
func (s *Service) EnforceUser(c echo.Context, ID int) error {
	// TODO: Implement querying db and checking the requested user's company_id
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

func (s *Service) isAdmin(c echo.Context) bool {
	return !(c.Get("role").(sandpiper.AccessRole) > sandpiper.AdminRole)
}

func (s *Service) isCompanyAdmin(c echo.Context) bool {
	// Must query company ID in database for the given user
	return !(c.Get("role").(sandpiper.AccessRole) > sandpiper.CompanyAdminRole)
}

// AccountCreate performs auth check when creating a new account
func (s *Service) AccountCreate(c echo.Context, roleID sandpiper.AccessRole, companyID uuid.UUID) error {
	if err := s.EnforceCompany(c, companyID); err != nil {
		return err
	}
	return s.IsLowerRole(c, roleID)
}

// IsLowerRole checks whether the requesting user has higher role than the user it wants to change
// Used for account creation/deletion
func (s *Service) IsLowerRole(c echo.Context, r sandpiper.AccessRole) error {
	return checkBool(c.Get("role").(sandpiper.AccessRole) < r)
}

func checkBool(b bool) error {
	if b {
		return nil
	}
	return echo.ErrForbidden
}

