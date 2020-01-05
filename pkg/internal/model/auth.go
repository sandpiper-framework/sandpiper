// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sandpiper

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// AuthToken holds authentication token details with refresh token
type AuthToken struct {
	Token        string `json:"token"`
	Expires      string `json:"expires"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshToken holds authentication token details
type RefreshToken struct {
	Token   string `json:"token"`
	Expires string `json:"expires"`
}

// AuthUser represents data ("claims") stored in JWT token for the current user
type AuthUser struct {
	ID        int
	CompanyID uuid.UUID
	Username  string
	Email     string
	Role      AccessLevel
}

// AtLeast checks user's access level
func (u *AuthUser) AtLeast(lvl AccessLevel) bool {
	// roles go from low to high (lower numbers are better)
	return u.Role <= lvl
}

// RBACService represents role-based access control service interface
type RBACService interface {
	CurrentUser(echo.Context) *AuthUser
	EnforceRole(echo.Context, AccessLevel) error
	EnforceUser(echo.Context, int) error
	EnforceCompany(echo.Context, uuid.UUID) error
	AccountCreate(echo.Context, AccessLevel, int, int) error
	IsLowerRole(echo.Context, AccessLevel) error
}