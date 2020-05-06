// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package activity contains services for sync history.
package activity

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"

	"sandpiper/pkg/api/activity/platform/pgsql"
	"sandpiper/pkg/shared/database"
	"sandpiper/pkg/shared/model"
)

// Service represents activity application interface (note no update!)
type Service interface {
	Create(echo.Context, sandpiper.Activity) (*sandpiper.Activity, error)
	List(echo.Context, *sandpiper.Pagination) ([]sandpiper.Activity, error)
	View(echo.Context, int) (*sandpiper.Activity, error)
	Delete(echo.Context, int) error
}

// New creates new activity application service
func New(db *database.DB, sdb Repository, rbac RBAC, sec Securer) *Activity {
	return &Activity{db: db.DB, sdb: sdb, rbac: rbac, sec: sec}
}

// Initialize initializes Activity application service with defaults
func Initialize(db *database.DB, rbac RBAC, sec Securer) *Activity {
	return New(db, pgsql.NewActivity(), rbac, sec)
}

// Activity represents activity application service
type Activity struct {
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
	Create(orm.DB, sandpiper.Activity) (*sandpiper.Activity, error)
	View(orm.DB, int) (*sandpiper.Activity, error)
	List(orm.DB, *sandpiper.Pagination) ([]sandpiper.Activity, error)
	Delete(orm.DB, int) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	EnforceRole(echo.Context, sandpiper.AccessLevel) error
}
