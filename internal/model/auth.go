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

// RBACService represents role-based access control service interface
type RBACService interface {
	CurrentUser(echo.Context) *AuthUser
	EnforceRole(echo.Context, AccessRole) error
	EnforceUser(echo.Context, int) error
	EnforceCompany(echo.Context, uuid.UUID) error
	AccountCreate(echo.Context, AccessRole, int, int) error
	IsLowerRole(echo.Context, AccessRole) error
}
