// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package slice

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/api/slice/platform/pgsql"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/database"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
)

// Service represents slice application interface
type Service interface {
	Create(echo.Context, sandpiper.Slice) (*sandpiper.Slice, error)
	List(echo.Context, *sandpiper.TagQuery, *sandpiper.Pagination) ([]sandpiper.Slice, error)
	View(echo.Context, uuid.UUID) (*sandpiper.Slice, error)
	ViewByName(echo.Context, string) (*sandpiper.Slice, error)
	Metadata(echo.Context, uuid.UUID) (sandpiper.MetaArray, error)
	Delete(echo.Context, uuid.UUID) error
	Update(echo.Context, *Update) (*sandpiper.Slice, error)
	Refresh(echo.Context, uuid.UUID) error
	Lock(echo.Context, uuid.UUID) error
	Unlock(echo.Context, uuid.UUID) error
}

// New creates new slice application service
func New(db *database.DB, sdb Repository, rbac RBAC, sec Securer) *Slice {
	return &Slice{db: db.DB, sdb: sdb, rbac: rbac, sec: sec}
}

// Initialize initializes Slice application service with defaults
func Initialize(db *database.DB, rbac RBAC, sec Securer) *Slice {
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
	ViewByName(orm.DB, uuid.UUID, string) (*sandpiper.Slice, error)
	List(orm.DB, *sandpiper.TagQuery, *sandpiper.Scope, *sandpiper.Pagination) ([]sandpiper.Slice, error)
	Metadata(orm.DB, uuid.UUID) (sandpiper.MetaArray, error)
	Update(orm.DB, *sandpiper.Slice) error
	Delete(orm.DB, *sandpiper.Slice) error
	Refresh(orm.DB, uuid.UUID) error
	Lock(orm.DB, uuid.UUID) error
	Unlock(orm.DB, uuid.UUID) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	CurrentUser(echo.Context) *sandpiper.AuthUser
	EnforceRole(echo.Context, sandpiper.AccessLevel) error
	EnforceScope(echo.Context) (*sandpiper.Scope, error)
}
