// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package handlers manages server-side rendering including signup, login and downloads
package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// LoginValues defines the login form fields that can be returned
type LoginValues struct {
	Email    string
	Password string
}

// Login handler
func Login(c echo.Context) error {

	// GET
	if c.Request().Method == http.MethodGet {
		vars := echo.Map{ // todo: replace this with config data
			"company": "Better Brakes",
			"terms":   "http://betterbrakes.com/terms",
		}
		return c.Render(http.StatusOK, "login.html", vars)
	}
	// POST
	result := new(LoginValues)
	if err := c.Bind(result); err != nil {
		return err
	}
	// todo: authenticate with loginValues, then save jwt in cookie if successful
	return c.Redirect(http.StatusFound, "/download")
}