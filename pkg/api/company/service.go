// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package company

// company service

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/params"

	"github.com/sandpiper-framework/sandpiper/pkg/api/company/platform/pgsql"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/database"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
)

// Service represents company application interface
type Service interface {
	Create(echo.Context, sandpiper.Company) (*sandpiper.Company, error)
	List(echo.Context, *params.Params) ([]sandpiper.Company, error)
	View(echo.Context, uuid.UUID) (*sandpiper.Company, error)
	Delete(echo.Context, uuid.UUID) error
	Update(echo.Context, *Update) (*sandpiper.Company, error)
	Server(echo.Context, uuid.UUID) (*sandpiper.Company, error)
	Servers(echo.Context, string) ([]sandpiper.Company, error)
}

// New creates new company application service
func New(db *database.DB, sdb Repository, rbac RBAC, sec Securer) *Company {
	return &Company{db: db.DB, sdb: sdb, rbac: rbac, sec: sec}
}

// Initialize initializes Company application service with defaults
func Initialize(db *database.DB, rbac RBAC, sec Securer) *Company {
	return New(db, pgsql.NewCompany(), rbac, sec)
}

// Company represents company application service
type Company struct {
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
	Create(orm.DB, sandpiper.Company) (*sandpiper.Company, error)
	View(orm.DB, uuid.UUID) (*sandpiper.Company, error)
	List(orm.DB, *sandpiper.Scope, *params.Params) ([]sandpiper.Company, error)
	Update(orm.DB, *sandpiper.Company) error
	Delete(orm.DB, *sandpiper.Company) error
	Server(orm.DB, uuid.UUID) (*sandpiper.Company, error)
	Servers(orm.DB, uuid.UUID, string) ([]sandpiper.Company, error)
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	CurrentUser(echo.Context) *sandpiper.AuthUser
	EnforceCompany(echo.Context, uuid.UUID) error
	EnforceRole(echo.Context, sandpiper.AccessLevel) error
	EnforceScope(echo.Context) (*sandpiper.Scope, error)
}
