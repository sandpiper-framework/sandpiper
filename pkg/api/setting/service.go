// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package setting

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"

	"sandpiper/pkg/api/setting/platform/pgsql"
	"sandpiper/pkg/shared/database"
	"sandpiper/pkg/shared/model"
)

// Service represents setting application interface (note no update!)
type Service interface {
	Create(echo.Context, *sandpiper.Setting) (*sandpiper.Setting, error)
	View(echo.Context) (*sandpiper.Setting, error)
	Update(echo.Context, *Update) (*sandpiper.Setting, error)
}

// New creates new setting application service
func New(db *database.DB, sdb Repository, rbac RBAC, sec Securer) *Setting {
	return &Setting{db: db.DB, sdb: sdb, rbac: rbac, sec: sec}
}

// Initialize initializes Setting application service with defaults
func Initialize(db *database.DB, rbac RBAC, sec Securer) *Setting {
	return New(db, pgsql.NewSetting(), rbac, sec)
}

// Setting represents setting application service
type Setting struct {
	db   *pg.DB
	sdb  Repository
	rbac RBAC
	sec  Securer
}

// Securer represents security interface
type Securer interface {
	Hash(string) string
}

// Repository represents available resource actions using a repository-abstraction-pattern interface.
type Repository interface {
	Create(orm.DB, *sandpiper.Setting) (*sandpiper.Setting, error)
	View(orm.DB) (*sandpiper.Setting, error)
	Update(orm.DB, *sandpiper.Setting) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	EnforceRole(echo.Context, sandpiper.AccessLevel) error
}
