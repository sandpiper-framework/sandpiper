// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package scope

import (
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/internal/model"
)

// Limit optionally adds a scope to list queries based on user roles
func Limit(u *sandpiper.AuthUser) (*sandpiper.Scoped, error) {
	switch true {
	case u.Role <= sandpiper.AdminRole: // user is SuperAdmin or Admin
		return nil, nil
	case u.Role == sandpiper.CompanyAdminRole:
		return &sandpiper.Scoped{Query: "company_id = ?", ID: u.CompanyID}, nil
	default:
		return nil, echo.ErrForbidden
	}
}
