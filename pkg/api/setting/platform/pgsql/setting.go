// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package pgsql

// setting service database access

import (
	"net/http"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
)

// Setting represents the client for setting table
type Setting struct{}

// NewSetting returns a new setting database instance
func NewSetting() *Setting {
	return &Setting{}
}

// Custom errors
var (
	// ErrSettingNotFound indicates select returned no rows
	ErrSettingNotFound = echo.NewHTTPError(http.StatusNotFound, "settings do not exist")
)

// Create creates a new setting in database (assumes allowed to do this).
func (s *Setting) Create(db orm.DB, setting *sandpiper.Setting) (*sandpiper.Setting, error) {
	if err := db.Insert(setting); err != nil {
		return nil, err
	}
	return setting, nil
}

// View returns a single setting by ID (assumes allowed to do this)
func (s *Setting) View(db orm.DB) (*sandpiper.Setting, error) {
	setting := &sandpiper.Setting{ID: true}
	err := db.Select(setting)
	if err != nil {
		return nil, selectError(err)
	}
	return setting, nil
}

// Update updates company info by primary key (assumes allowed to do this)
func (s *Setting) Update(db orm.DB, setting *sandpiper.Setting) error {
	return db.Update(setting)
}

func selectError(err error) error {
	if err == pg.ErrNoRows {
		return ErrSettingNotFound
	}
	return err
}
