package mock

import (
	"github.com/labstack/echo/v4"
	"github.com/satori/go.uuid"

	"autocare.org/sandpiper/internal/model"
)

// RBAC Mock
type RBAC struct {
	UserFn            func(echo.Context) *sandpiper.AuthUser
	EnforceRoleFn     func(echo.Context, sandpiper.AccessRole) error
	EnforceUserFn     func(echo.Context, int) error
	EnforceCompanyFn  func(echo.Context, uuid.UUID) error
	AccountCreateFn   func(echo.Context, sandpiper.AccessRole, uuid.UUID) error
	IsLowerRoleFn     func(echo.Context, sandpiper.AccessRole) error
}

// User mock
func (a *RBAC) User(c echo.Context) *sandpiper.AuthUser {
	return a.UserFn(c)
}

// EnforceRole mock
func (a *RBAC) EnforceRole(c echo.Context, role sandpiper.AccessRole) error {
	return a.EnforceRoleFn(c, role)
}

// EnforceUser mock
func (a *RBAC) EnforceUser(c echo.Context, id int) error {
	return a.EnforceUserFn(c, id)
}

// EnforceCompany mock
func (a *RBAC) EnforceCompany(c echo.Context, id uuid.UUID) error {
	return a.EnforceCompanyFn(c, id)
}

// AccountCreate mock
func (a *RBAC) AccountCreate(c echo.Context, roleID sandpiper.AccessRole, companyID uuid.UUID) error {
	return a.AccountCreateFn(c, roleID, companyID)
}

// IsLowerRole mock
func (a *RBAC) IsLowerRole(c echo.Context, role sandpiper.AccessRole) error {
	return a.IsLowerRoleFn(c, role)
}
