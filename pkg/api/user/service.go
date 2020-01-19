// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package user

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/user/platform/pgsql"
	"autocare.org/sandpiper/pkg/shared/model"
)

// Service represents user application interface
type Service interface {
	Create(echo.Context, sandpiper.User) (*sandpiper.User, error)
	List(echo.Context, *sandpiper.Pagination) ([]sandpiper.User, error)
	View(echo.Context, int) (*sandpiper.User, error)
	Delete(echo.Context, int) error
	Update(echo.Context, *Update) (*sandpiper.User, error)
}

// User represents user application service
type User struct {
	db   *pg.DB
	sdb  Repository
	rbac RBAC
	sec  Securer
}

// New creates new user application service
func New(db *pg.DB, sdb Repository, rbac RBAC, sec Securer) *User {
	return &User{db: db, sdb: sdb, rbac: rbac, sec: sec}
}

// Initialize initializes User application service with defaults
func Initialize(db *pg.DB, rbac RBAC, sec Securer) *User {
	return New(db, pgsql.NewUser(), rbac, sec)
}

// Securer represents security interface
type Securer interface {
	Hash(string) string
}

// Repository represents available resource actions using a repository-abstraction-pattern interface.
type Repository interface {
	Create(orm.DB, sandpiper.User) (*sandpiper.User, error)
	View(orm.DB, int) (*sandpiper.User, error)
	List(orm.DB, *sandpiper.Clause, *sandpiper.Pagination) ([]sandpiper.User, error)
	Update(orm.DB, *sandpiper.User) error
	Delete(orm.DB, *sandpiper.User) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	CurrentUser(echo.Context) *sandpiper.AuthUser
	EnforceUser(echo.Context, int) error
	AccountCreate(echo.Context, sandpiper.AccessLevel, uuid.UUID) error
	IsLowerRole(echo.Context, sandpiper.AccessLevel) error
	EnforceScope(echo.Context) (*sandpiper.Clause, error)
}
