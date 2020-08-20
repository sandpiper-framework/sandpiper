// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package web manages server-side rendering including signup, login and downloads
package web

import (
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/echoview-v4"
	"github.com/foolin/goview/supports/gorice"
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/api/web/handlers"
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
	srv.GET("/", handlers.Login)
	srv.POST("/", handlers.Login)

	// signup page
	srv.GET("/signup", handlers.Signup)
	srv.POST("/signup", handlers.Signup)

	// download page
	srv.GET("/download", handlers.Download)
	srv.POST("/download", handlers.Download)
}
