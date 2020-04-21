// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sync

import (
	"autocare.org/sandpiper/pkg/shared/database"
	"net/url"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/sync/platform/pgsql"
	"autocare.org/sandpiper/pkg/shared/model"
)

// Service represents sync application interface
type Service interface {
	Start(echo.Context, *url.URL) error
	Connect(echo.Context) error
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
	LogActivity(orm.DB, sandpiper.SyncRequest) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	CurrentUser(echo.Context) *sandpiper.AuthUser
	EnforceRole(echo.Context, sandpiper.AccessLevel) error
	EnforceScope(echo.Context) (*sandpiper.Scope, error)
	EnforceServerRole(string) error
}
