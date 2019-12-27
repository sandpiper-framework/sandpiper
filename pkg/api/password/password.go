// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package password manages the password service which allows setting a new password,
// validating password strength, and generating the encrypted hash to save in the db.
package password

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Custom errors
var (
	ErrIncorrectPassword = echo.NewHTTPError(http.StatusBadRequest, "incorrect old password")
	ErrInsecurePassword  = echo.NewHTTPError(http.StatusBadRequest, "insecure password")
)

// Change changes user's password
func (p *Password) Change(c echo.Context, userID int, oldPass, newPass string) error {
	if err := p.rbac.EnforceUser(c, userID); err != nil {
		return err
	}

	u, err := p.udb.View(p.db, userID)
	if err != nil {
		return err
	}

	if !p.sec.HashMatchesPassword(u.Password, oldPass) {
		return ErrIncorrectPassword
	}

	if !p.sec.Password(newPass, u.FirstName, u.LastName, u.Username, u.Email) {
		return ErrInsecurePassword
	}

	u.ChangePassword(p.sec.Hash(newPass))

	return p.udb.Update(p.db, u)
}
