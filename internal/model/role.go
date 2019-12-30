// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sandpiper

import (
	"github.com/google/uuid"
)

// AccessRole represents access role type
type AccessRole int

const (
	// SuperAdminRole has all permissions
	SuperAdminRole AccessRole = 100

	// AdminRole has admin specific permissions
	AdminRole AccessRole = 110

	// CompanyAdminRole can edit company specific things
	CompanyAdminRole AccessRole = 120

	// UserRole is a standard user
	UserRole AccessRole = 200
)

// Scoped adds additional restrictions for scoping list queries based on roles
type Scoped struct {
	Query string
	ID    uuid.UUID  // usually a companyID
}

