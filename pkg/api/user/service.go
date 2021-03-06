// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package user

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/params"

	"github.com/sandpiper-framework/sandpiper/pkg/api/user/platform/pgsql"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/database"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
)

// Service represents user application interface
type Service interface {
	Create(echo.Context, sandpiper.User) (*sandpiper.User, error)
	List(echo.Context, *params.Params) ([]sandpiper.User, error)
	View(echo.Context, int) (*sandpiper.User, error)
	Delete(echo.Context, int) error
	Update(echo.Context, *Update) (*sandpiper.User, error)
	CreateAPIKey(echo.Context) (*sandpiper.APIKey, error)
}

// User represents user application service
type User struct {
	db   *pg.DB
	sdb  Repository
	rbac RBAC
	sec  Securer
}

// New creates new user application service
func New(db *database.DB, sdb Repository, rbac RBAC, sec Securer) *User {
	return &User{db: db.DB, sdb: sdb, rbac: rbac, sec: sec}
}

// Initialize initializes User application service with defaults
func Initialize(db *database.DB, rbac RBAC, sec Securer) *User {
	return New(db, pgsql.NewUser(), rbac, sec)
}

// Securer represents security interface
type Securer interface {
	Hash(string) string
	APIKeySecret() string
	RandomPassword(int) (string, error)
}

// Repository represents available resource actions using a repository-abstraction-pattern interface.
type Repository interface {
	Create(orm.DB, sandpiper.User) (*sandpiper.User, error)
	View(orm.DB, int) (*sandpiper.User, error)
	List(orm.DB, *params.Params, *sandpiper.Scope) ([]sandpiper.User, error)
	Update(orm.DB, *sandpiper.User) error
	Delete(orm.DB, *sandpiper.User) error
	CompanySyncUser(orm.DB, uuid.UUID) (*sandpiper.User, error)
	UpdateSyncUser(orm.DB, *sandpiper.User) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	CurrentUser(echo.Context) *sandpiper.AuthUser
	OurServer() *sandpiper.Server
	EnforceUser(echo.Context, int) error
	AccountCreate(echo.Context, sandpiper.AccessLevel, uuid.UUID) error
	IsLowerRole(echo.Context, sandpiper.AccessLevel) error
	EnforceScope(echo.Context) (*sandpiper.Scope, error)
	EnforceRole(echo.Context, sandpiper.AccessLevel) error
	EnforceServerRole(string) error
}
