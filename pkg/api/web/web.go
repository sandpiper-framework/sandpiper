// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package web manages server-side rendering including signup, login and downloads
package web

import (
	"net/http"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/echoview-v4"
	"github.com/foolin/goview/supports/gorice"
	"github.com/labstack/echo/v4"
)

// FileServer serves static files and templates from embedded files in `rice-box.go`
func FileServer(srv *echo.Echo) {
	// set view engine (for templates)
	viewConfig := goview.DefaultConfig
	// viewConfig.DisableCache = true // auto reload template file for debug.
	views := gorice.NewWithConfig(rice.MustFindBox("views"), viewConfig)
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

	// download page
	srv.GET("/download", download)
	srv.POST("/download", download)
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
	// NOTE: I'm not sure this is a proper use of redirect. Not sure how to do it another way.
	return c.Redirect(http.StatusFound, "/download")
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
	// todo: do something with signupValues (maybe email or save to table)
	// display an Acknowledgment
	vars["thankyou"] = true
	return c.Render(http.StatusOK, "signup.html", vars)
}

func download(c echo.Context) error {
	type tableRow struct {
		SliceName   string
		SliceType   string
		ContentDate string
	}
	type tableData struct {
		Downloading bool
		Rows []tableRow
	}

	// GET
	if c.Request().Method == http.MethodGet {
		// todo: pull from database
		data := tableData{
			Rows: []tableRow{
				{"slice1", "aces-file", time.Now().Format("2006-01-02")},
				{"slice2", "pies-file", time.Now().Format("2006-01-02")},
			},
		}
		return c.Render(http.StatusOK, "download.html", data)
	}

	// POST
	vars := echo.Map{"Downloading": true}
	return c.Render(http.StatusOK, "download.html", vars)
}
