// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package web manages server-side rendering including signup and login
package web

import (
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/foolin/goview/supports/echoview-v4"
	"github.com/foolin/goview/supports/gorice"
	"github.com/labstack/echo/v4"
)

type signupValues struct {
	Name string
	Email string
	Company string
	SandpiperID string
	Kind string
}

// FileServer serves static files and templates from embedded files in `rice-box.go`
func FileServer(srv *echo.Echo) {
	// static files file server
	box := rice.MustFindBox("static")
	static := http.StripPrefix("/static/", http.FileServer(box.HTTPBox()))

	// set view engine (for templates)
	views := gorice.New(rice.MustFindBox("views"))
	srv.Renderer = echoview.Wrap(views)

	// routes
	srv.GET("/", login)
	srv.GET("/signup", signup)
	srv.GET("/static/*", echo.WrapHandler(static))
}

func login(c echo.Context) error {
	// full name with extension to render just the file
	vars := echo.Map{	// todo: replace this with config data
		"company": "Better Brakes",
		"terms": "http://betterbrakes.com/terms",
	}
	return c.Render(http.StatusOK, "login.html", vars)
}

func signup(c echo.Context) error {
	// full name with extension to render just the file
	vars := echo.Map{	// todo: replace this with config data
		"company": "Better Brakes",
		"terms": "http://betterbrakes.com/terms",
	}
	return c.Render(http.StatusOK, "signup.html", vars)
}
