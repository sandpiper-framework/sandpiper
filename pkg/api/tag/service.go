// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package tag

// tag service

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/api/tag/platform/pgsql"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/database"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
)

// Service represents tag application interface
type Service interface {
	Create(echo.Context, sandpiper.Tag) (*sandpiper.Tag, error)
	List(echo.Context, *sandpiper.Pagination) ([]sandpiper.Tag, error)
	View(echo.Context, int) (*sandpiper.Tag, error)
	Delete(echo.Context, int) error
	Update(echo.Context, *Update) (*sandpiper.Tag, error)
	Assign(echo.Context, int, uuid.UUID) error
	Remove(echo.Context, int, uuid.UUID) error
}

// New creates new company application service
func New(db *database.DB, sdb Repository, rbac RBAC, sec Securer) *Tag {
	return &Tag{db: db.DB, sdb: sdb, rbac: rbac, sec: sec}
}

// Initialize initializes Tag application service with defaults
func Initialize(db *database.DB, rbac RBAC, sec Securer) *Tag {
	return New(db, pgsql.NewTag(), rbac, sec)
}

// Tag represents company application service
type Tag struct {
	db   *pg.DB
	sdb  Repository // service repository database interface
	rbac RBAC
	sec  Securer
}

// Securer represents security interface
type Securer interface {
	Hash(string) string
}

// Repository represents available resource actions using a repository-abstraction-pattern interface.
type Repository interface {
	Create(orm.DB, sandpiper.Tag) (*sandpiper.Tag, error)
	View(orm.DB, int) (*sandpiper.Tag, error)
	List(orm.DB, *sandpiper.Pagination) ([]sandpiper.Tag, error)
	Update(orm.DB, *sandpiper.Tag) error
	Delete(orm.DB, *sandpiper.Tag) error
	Assign(orm.DB, int, uuid.UUID) error
	Remove(orm.DB, int, uuid.UUID) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	CurrentUser(echo.Context) *sandpiper.AuthUser
	EnforceRole(echo.Context, sandpiper.AccessLevel) error
}
