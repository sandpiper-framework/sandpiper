package query

import (
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/model"
)

// List prepares data for list queries
func List(u *sandpiper.AuthUser) (*sandpiper.ListQuery, error) {
	switch true {
	case u.Role <= sandpiper.AdminRole: // user is SuperAdmin or Admin
		return nil, nil
	case u.Role == sandpiper.CompanyAdminRole:
		return &sandpiper.ListQuery{Query: "company_id = ?", ID: u.CompanyID}, nil
	case u.Role == sandpiper.LocationAdminRole:
		return &sandpiper.ListQuery{Query: "location_id = ?", ID: u.LocationID}, nil
	default:
		return nil, echo.ErrForbidden
	}
}
