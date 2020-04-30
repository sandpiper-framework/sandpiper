// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sync

import (
	"autocare.org/sandpiper/pkg/shared/database"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/sync/platform/pgsql"
	"autocare.org/sandpiper/pkg/shared/model"
)

// Service represents sync application interface
type Service interface {
	Start(echo.Context, uuid.UUID) error
	Process(echo.Context) error
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
	key  string
}

// Securer represents security interface
type Securer interface {
	Hash(string) string
	APIKeySecret() string
}

// Repository represents available resource actions using a repository-abstraction-pattern interface.
type Repository interface {
	Primary(orm.DB, uuid.UUID) (*sandpiper.Company, error)
	LogActivity(orm.DB, sandpiper.SyncRequest) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	CurrentUser(echo.Context) *sandpiper.AuthUser
	EnforceRole(echo.Context, sandpiper.AccessLevel) error
	EnforceScope(echo.Context) (*sandpiper.Scope, error)
	EnforceServerRole(string) error
}
