// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package query

import (
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/internal/model"
)

// List prepares data for list queries
func List(u *sandpiper.AuthUser) (*sandpiper.ListQuery, error) {
	switch true {
	case u.Role <= sandpiper.AdminRole: // user is SuperAdmin or Admin
		return nil, nil
	case u.Role == sandpiper.CompanyAdminRole:
		return &sandpiper.ListQuery{Query: "company_id = ?", ID: u.CompanyID}, nil
	default:
		return nil, echo.ErrForbidden
	}
}
