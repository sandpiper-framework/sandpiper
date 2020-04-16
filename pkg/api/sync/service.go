// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package sync contains services for sync history.
package sync

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/sync/platform/pgsql"
	"autocare.org/sandpiper/pkg/shared/database"
	"autocare.org/sandpiper/pkg/shared/model"
)

// Service represents sync application interface (note no update!)
type Service interface {
	Create(echo.Context, sandpiper.Sync) (*sandpiper.Sync, error)
	List(echo.Context, *sandpiper.Pagination) ([]sandpiper.Sync, error)
	View(echo.Context, int) (*sandpiper.Sync, error)
	Delete(echo.Context, int) error
}

// New creates new sync application service
func New(db *database.DB, sdb Repository, rbac RBAC, sec Securer) *Sync {
	return &Sync{db: db.DB, sdb: sdb, rbac: rbac, sec: sec}
}

// Initialize initializes Sync application service with defaults
func Initialize(db *database.DB, rbac RBAC, sec Securer) *Sync {
	return New(db, pgsql.NewSync(), rbac, sec)
}

// Sync represents sync application service
type Sync struct {
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
	Create(orm.DB, sandpiper.Sync) (*sandpiper.Sync, error)
	//CompanySubscribed(db orm.DB, companyID uuid.UUID, syncID int) bool
	View(orm.DB, int) (*sandpiper.Sync, error)
	List(orm.DB, *sandpiper.Scope, *sandpiper.Pagination) ([]sandpiper.Sync, error)
	Delete(orm.DB, int) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	CurrentUser(echo.Context) *sandpiper.AuthUser
	EnforceRole(echo.Context, sandpiper.AccessLevel) error
	EnforceScope(echo.Context) (*sandpiper.Scope, error)
}
