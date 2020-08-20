// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package handlers manages server-side rendering including signup, login and downloads
package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// SignupValues defines the signup form fields that can be returned
type SignupValues struct {
	Name     string
	Email    string
	Company  string
	ServerID string
	Kind     string
}

// Signup handler
func Signup(c echo.Context) error {

	vars := echo.Map{ // todo: replace this with config data
		"company": "Better Brakes",
		"terms":   "http://betterbrakes.com/terms",
	}

	// GET
	if c.Request().Method == http.MethodGet {
		// render signup page
		return c.Render(http.StatusOK, "signup.html", vars)
	}

	// POST
	result := new(SignupValues)
	if err := c.Bind(result); err != nil {
		return err
	}
	// todo: do something with signupValues (maybe email or save to table)
	// display an Acknowledgment
	vars["thankyou"] = true
	return c.Render(http.StatusOK, "signup.html", vars)
}
