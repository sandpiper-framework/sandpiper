// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package web manages server-side rendering including signup, login and downloads
package web

import (
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/foolin/goview/supports/echoview-v4"
	"github.com/foolin/goview/supports/gorice"
	"github.com/labstack/echo/v4"
)

// FileServer serves static files and templates from embedded files in `rice-box.go`
func FileServer(srv *echo.Echo) {
	// set view engine (for templates)
	views := gorice.New(rice.MustFindBox("views"))
	srv.Renderer = echoview.Wrap(views)

	// handle all static assets
	box := rice.MustFindBox("static")
	static := http.StripPrefix("/static/", http.FileServer(box.HTTPBox()))
	srv.GET("/static/*", echo.WrapHandler(static))

	// login page
	srv.GET("/", login)
	srv.POST("/", login)

	// signup page
	srv.GET("/signup", signup)
	srv.POST("/signup", signup)

	// todo: download page
	// srv.GET("/download", download)
	// srv.POST("/download", download)
}

func login(c echo.Context) error {
	type loginValues struct {
		Email    string
		Password string
	}
	// GET
	if c.Request().Method == http.MethodGet {
		vars := echo.Map{ // todo: replace this with config data
			"company": "Better Brakes",
			"terms":   "http://betterbrakes.com/terms",
		}
		return c.Render(http.StatusOK, "login.html", vars)
	}
	// POST
	result := new(loginValues)
	if err := c.Bind(result); err != nil {
		return err
	}
	// todo: authenticate with loginValues, then save jwt in cookie and redirect to download.
	return c.JSON(http.StatusOK, result)
}

func signup(c echo.Context) error {
	type signupValues struct {
		Name     string
		Email    string
		Company  string
		ServerID string
		Kind     string
	}

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
	result := new(signupValues)
	if err := c.Bind(result); err != nil {
		return err
	}
	// todo: do something with signupValues (maybew email or save to table)
	// display an Acknowledgment
	vars["thankyou"] = true
	return c.Render(http.StatusOK, "signup.html", vars)
}
