// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package grain

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"sandpiper/pkg/api/grain/platform/pgsql"
	"sandpiper/pkg/shared/database"
	"sandpiper/pkg/shared/model"
)

// Service represents grain application interface (note no update!)
type Service interface {
	Create(echo.Context, bool, *sandpiper.Grain) (*sandpiper.Grain, error)
	List(echo.Context, bool, *sandpiper.Pagination) ([]sandpiper.Grain, error)
	ListBySlice(echo.Context, uuid.UUID, bool, *sandpiper.Pagination) ([]sandpiper.Grain, error)
	View(echo.Context, uuid.UUID) (*sandpiper.Grain, error)
	ViewByKeys(echo.Context, uuid.UUID, string, bool) (*sandpiper.Grain, error)
	Delete(echo.Context, uuid.UUID) error
}

// New creates new grain application service
func New(db *database.DB, sdb Repository, rbac RBAC, sec Securer) *Grain {
	return &Grain{db: db.DB, sdb: sdb, rbac: rbac, sec: sec}
}

// Initialize initializes Grain application service with defaults
func Initialize(db *database.DB, rbac RBAC, sec Securer) *Grain {
	return New(db, pgsql.NewGrain(), rbac, sec)
}

// Grain represents grain application service
type Grain struct {
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
	Create(orm.DB, bool, *sandpiper.Grain) (*sandpiper.Grain, error)
	CompanySubscribed(db orm.DB, companyID uuid.UUID, grainID uuid.UUID) bool
	View(orm.DB, uuid.UUID) (*sandpiper.Grain, error)
	ViewByKeys(orm.DB, uuid.UUID, string, bool) (*sandpiper.Grain, error)
	List(orm.DB, uuid.UUID, bool, *sandpiper.Scope, *sandpiper.Pagination) ([]sandpiper.Grain, error)
	Delete(orm.DB, uuid.UUID) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	CurrentUser(echo.Context) *sandpiper.AuthUser
	EnforceRole(echo.Context, sandpiper.AccessLevel) error
	EnforceScope(echo.Context) (*sandpiper.Scope, error)
}
