// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package auth is for identity and access management
package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"sandpiper/pkg/shared/model"
)

// Custom errors
var (
	ErrInvalidCredentials = echo.NewHTTPError(http.StatusUnauthorized, "Authentication Error.")
	ErrNotAuthorized      = echo.NewHTTPError(http.StatusUnauthorized, "User is not authorized")
)

// Authenticate tries to authenticate the user provided by username and password
func (a *Auth) Authenticate(c echo.Context, user, pass string) (*sandpiper.AuthToken, error) {
	u, err := a.sdb.FindByUsername(a.db, user)
	if err != nil {
		return nil, err
	}

	if !a.sec.HashMatchesPassword(u.Password, pass) {
		return nil, ErrInvalidCredentials
	}

	if !u.Active {
		return nil, ErrNotAuthorized
	}

	// generate new jwt with user information in the claims
	token, expire, err := a.tg.GenerateToken(u)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	// save token and when last login in user record
	u.UpdateLastLogin(a.sec.Token(token))

	if err := a.sdb.Update(a.db, u); err != nil {
		return nil, err
	}

	return &sandpiper.AuthToken{Token: token, Expires: expire, RefreshToken: u.Token}, nil
}

// Refresh refreshes jwt token and puts new claims inside
func (a *Auth) Refresh(c echo.Context, refreshToken string) (*sandpiper.RefreshToken, error) {
	user, err := a.sdb.FindByToken(a.db, refreshToken)
	if err != nil {
		return nil, err
	}
	if !user.Active {
		return nil, ErrNotAuthorized
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
	return a.sdb.View(a.db, au.ID)
}
