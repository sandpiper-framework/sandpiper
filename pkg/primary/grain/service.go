// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package grain

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/primary/grain/platform/pgsql"
	"autocare.org/sandpiper/pkg/shared/model"
)

// Service represents grain application interface (note no update!)
type Service interface {
	Create(echo.Context, bool, sandpiper.Grain) (*sandpiper.Grain, error)
	List(echo.Context, bool, *sandpiper.Pagination) ([]sandpiper.Grain, error)
	View(echo.Context, uuid.UUID) (*sandpiper.Grain, error)
	Exists(echo.Context, uuid.UUID, string, string) (*sandpiper.Grain, error)
	Delete(echo.Context, uuid.UUID) error
}

// New creates new grain application service
func New(db *pg.DB, sdb Repository, rbac RBAC, sec Securer) *Grain {
	return &Grain{db: db, sdb: sdb, rbac: rbac, sec: sec}
}

// Initialize initializes Grain application service with defaults
func Initialize(db *pg.DB, rbac RBAC, sec Securer) *Grain {
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
	Create(orm.DB, bool, sandpiper.Grain) (*sandpiper.Grain, error)
	CompanySubscribed(db orm.DB, companyID uuid.UUID, grainID uuid.UUID) bool
	View(orm.DB, uuid.UUID) (*sandpiper.Grain, error)
	Exists(orm.DB, uuid.UUID, string, string) (*sandpiper.Grain, error)
	List(orm.DB, bool, *sandpiper.Scope, *sandpiper.Pagination) ([]sandpiper.Grain, error)
	Delete(orm.DB, uuid.UUID) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	CurrentUser(echo.Context) *sandpiper.AuthUser
	EnforceRole(echo.Context, sandpiper.AccessLevel) error
	EnforceScope(echo.Context) (*sandpiper.Scope, error)
}
