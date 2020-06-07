// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package auth is for identity and access management
package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/secure"
)

// Custom errors
var (
	ErrInvalidCredentials = echo.NewHTTPError(http.StatusUnauthorized, "Authentication Error.")
	ErrNotAuthorized      = echo.NewHTTPError(http.StatusUnauthorized, "User is not authorized")
	ErrMissingCredentials = echo.NewHTTPError(http.StatusBadRequest, "missing required credentials")
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

// Server returns info about the server
func (a *Auth) Server(c echo.Context) *sandpiper.Server {
	return a.rbac.OurServer()
}

// ParseCredentials extracts username/password from a request body (in plain text). If
// api-key is supplied, it decrypts the api-key into those plain text values instead
func (a *Auth) ParseCredentials(c echo.Context) (*secure.Credentials, error) {
	var err error

	creds := new(secure.Credentials)
	if err := c.Bind(creds); err != nil {
		return nil, err
	}
	if creds.SyncAPIKey != "" {
		// extract username and password from api-key
		creds, err = secure.NewCredentials(creds.SyncAPIKey, a.sec.APIKeySecret())
		if err != nil {
			return nil, err
		}
	}
	if creds.Username == "" || creds.Password == "" {
		return nil, ErrMissingCredentials
	}
	return creds, nil
}
