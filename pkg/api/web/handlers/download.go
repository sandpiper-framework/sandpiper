// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package handlers manages server-side rendering including signup, login and downloads
package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// Download handler
func Download(c echo.Context) error {
	type tableRow struct {
		SliceName   string
		SliceType   string
		ContentDate string
	}
	type tableData struct {
		Downloading bool
		Rows        []tableRow
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
