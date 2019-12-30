// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package scope

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/internal/model"
)

// Clause adds additional restrictions for scoping list queries based on roles
type Clause struct {
	Condition string
	ID    uuid.UUID  // usually a companyID
}

// Limit optionally adds a scope to list queries based on user roles
func Limit(u *sandpiper.AuthUser) (*Clause, error) {
	switch true {
	case u.Role <= sandpiper.AdminRole: // user is SuperAdmin or Admin
		return nil, nil
	case u.Role == sandpiper.CompanyAdminRole:
		return &Clause{Condition: "company_id = ?", ID: u.CompanyID}, nil
	default:
		return nil, echo.ErrForbidden
	}
}
