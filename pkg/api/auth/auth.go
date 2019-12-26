// Package auth is for identity and access management
package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/internal/model"
)

// Custom errors
var (
	ErrInvalidCredentials = echo.NewHTTPError(http.StatusUnauthorized, "Username or password does not exist")
	ErrNotAuthorized = echo.NewHTTPError(http.StatusUnauthorized, "User is not authorized")
)

// Authenticate tries to authenticate the user provided by username and password
func (a *Auth) Authenticate(c echo.Context, user, pass string) (*sandpiper.AuthToken, error) {
	u, err := a.udb.FindByUsername(a.db, user)
	if err != nil {
		return nil, err
	}

	if !a.sec.HashMatchesPassword(u.Password, pass) {
		return nil, ErrInvalidCredentials
	}

	if !u.Active {
		return nil, ErrNotAuthorized
	}

	token, expire, err := a.tg.GenerateToken(u)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	u.UpdateLastLogin(a.sec.Token(token))

	if err := a.udb.Update(a.db, u); err != nil {
		return nil, err
	}

	return &sandpiper.AuthToken{Token: token, Expires: expire, RefreshToken: u.Token}, nil
}

// Refresh refreshes jwt token and puts new claims inside
func (a *Auth) Refresh(c echo.Context, token string) (*sandpiper.RefreshToken, error) {
	user, err := a.udb.FindByToken(a.db, token)
	if err != nil {
		return nil, err
	}
	token, expire, err := a.tg.GenerateToken(user)
	if err != nil {
		return nil, err
	}
	return &sandpiper.RefreshToken{Token: token, Expires: expire}, nil
}

// Me returns info about currently logged in user
func (a *Auth) Me(c echo.Context) (*sandpiper.User, error) {
	au := a.rbac.CurrentUser(c)
	return a.udb.View(a.db, au.ID)
}
