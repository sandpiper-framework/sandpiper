// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sandpiper

// AccessLevel represents user access role (in a hierarchy of levels)
type AccessLevel int

const (
	// SuperAdminRole has all permissions
	SuperAdminRole AccessLevel = 100

	// AdminRole has admin specific permissions
	AdminRole AccessLevel = 110

	// CompanyAdminRole can edit company specific things
	CompanyAdminRole AccessLevel = 120

	// UserRole is a standard user
	UserRole AccessLevel = 200
)
