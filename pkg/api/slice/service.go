// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package slice

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/internal/model"
	"autocare.org/sandpiper/pkg/internal/scope"
	"autocare.org/sandpiper/pkg/api/slice/platform/pgsql"
)

// Service represents slice application interface
type Service interface {
	Create(echo.Context, sandpiper.Slice) (*sandpiper.Slice, error)
	List(echo.Context, *sandpiper.Pagination) ([]sandpiper.Slice, error)
	View(echo.Context, uuid.UUID) (*sandpiper.Slice, error)
	Delete(echo.Context, uuid.UUID) error
	Update(echo.Context, *Update) (*sandpiper.Slice, error)
}

// New creates new slice application service
func New(db *pg.DB, sdb Repository, rbac RBAC, sec Securer) *Slice {
	return &Slice{db: db, sdb: sdb, rbac: rbac, sec: sec}
}

// Initialize initializes Slice application service with defaults
func Initialize(db *pg.DB, rbac RBAC, sec Securer) *Slice {
	return New(db, pgsql.NewSlice(), rbac, sec)
}

// Slice represents slice application service
type Slice struct {
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
	Create(orm.DB, sandpiper.Slice) (*sandpiper.Slice, error)
	View(orm.DB, uuid.UUID) (*sandpiper.Slice, error)
	ViewBySub(db orm.DB, companyID uuid.UUID, sliceID uuid.UUID) (*sandpiper.Slice, error)
	List(orm.DB, *scope.Clause, *sandpiper.Pagination) ([]sandpiper.Slice, error)
	Update(orm.DB, *sandpiper.Slice) error
	Delete(orm.DB, *sandpiper.Slice) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	CurrentUser(echo.Context) *sandpiper.AuthUser
	EnforceRole(echo.Context, sandpiper.AccessLevel) error
}
