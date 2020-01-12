// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package grain

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/grain/platform/pgsql"
	"autocare.org/sandpiper/pkg/internal/model"
	"autocare.org/sandpiper/pkg/internal/scope"
)

// Service represents grain application interface
type Service interface {
	Create(echo.Context, sandpiper.Grain) (*sandpiper.Grain, error)
	List(echo.Context, *sandpiper.Pagination) ([]sandpiper.Grain, error)
	View(echo.Context, uuid.UUID) (*sandpiper.Grain, error)
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
	Create(orm.DB, sandpiper.Grain) (*sandpiper.Grain, error)
	View(orm.DB, uuid.UUID) (*sandpiper.Grain, error)
	ViewBySlice(orm.DB, uuid.UUID) (*sandpiper.Grain, error)
	ViewBySub(db orm.DB, companyID uuid.UUID, sliceID uuid.UUID) (*sandpiper.Grain, error)
	List(orm.DB, *scope.Clause, *sandpiper.Pagination) ([]sandpiper.Grain, error)
	Delete(orm.DB, *sandpiper.Grain) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	CurrentUser(echo.Context) *sandpiper.AuthUser
	EnforceRole(echo.Context, sandpiper.AccessLevel) error
}
