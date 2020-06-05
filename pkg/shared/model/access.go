// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package sandpiper

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// AccessLevel represents user access role (in a hierarchy of levels)
type AccessLevel int

const (
	// SuperAdminRole has all permissions
	SuperAdminRole AccessLevel = 100

	// AdminRole has admin specific permissions
	AdminRole AccessLevel = 110

	// CompanyAdminRole can maintain company-specific things
	CompanyAdminRole AccessLevel = 120

	// SyncRole is a special limited-access user for the sync process
	SyncRole AccessLevel = 200
)

// RoleIsValid validates against available access levels
func RoleIsValid(role AccessLevel) bool {
	switch role {
	case SuperAdminRole, AdminRole, CompanyAdminRole, SyncRole:
		return true
	default:
		return false
	}
}

// AtLeast checks user's access level
func (u *AuthUser) AtLeast(lvl AccessLevel) bool {
	// roles go from low to high (lower numbers are better)
	return u.Role <= lvl
}

// Scope adds additional restrictions for scoping list queries based on roles
type Scope struct {
	Condition string
	ID        uuid.UUID // always a companyID
}

// ApplyScope adds any needed scope to a list-query based on user roles
func (u *AuthUser) ApplyScope(lhs string) (*Scope, error) {
	switch true {
	case u.Role <= AdminRole: // user is SuperAdmin or Admin
		return nil, nil
	case u.Role == CompanyAdminRole:
		return &Scope{Condition: lhs + " = ?", ID: u.CompanyID}, nil
	default: // is standard user
		return nil, echo.ErrForbidden
	}
}
