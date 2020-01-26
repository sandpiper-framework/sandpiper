// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package tag

// tag service

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/tag/platform/pgsql"
	"autocare.org/sandpiper/pkg/shared/model"
)

// Service represents tag application interface
type Service interface {
	Create(echo.Context, sandpiper.Tag) (*sandpiper.Tag, error)
	List(echo.Context, *sandpiper.Pagination) ([]sandpiper.Tag, error)
	View(echo.Context, int) (*sandpiper.Tag, error)
	Delete(echo.Context, int) error
	Update(echo.Context, *Update) (*sandpiper.Tag, error)
}

// New creates new company application service
func New(db *pg.DB, sdb Repository, rbac RBAC, sec Securer) *Tag {
	return &Tag{db: db, sdb: sdb, rbac: rbac, sec: sec}
}

// Initialize initializes Tag application service with defaults
func Initialize(db *pg.DB, rbac RBAC, sec Securer) *Tag {
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
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	CurrentUser(echo.Context) *sandpiper.AuthUser
	EnforceRole(echo.Context, sandpiper.AccessLevel) error
}
