// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package static creates routes for serving static files
package static

import (
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/labstack/echo/v4"
)

// FileServer serves from embedded files in `rice-box.go`
func FileServer(srv *echo.Echo) {
	// relative to this source file (three levels down!)
	public := http.FileServer(rice.MustFindBox("../../../public").HTTPBox())
	srv.GET("/", echo.WrapHandler(public))
	srv.GET("/*/*", echo.WrapHandler(http.StripPrefix("/", public)))
}
