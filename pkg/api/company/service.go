// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package company

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/internal/model"
	"autocare.org/sandpiper/pkg/internal/scope"
	"autocare.org/sandpiper/pkg/api/company/platform/pgsql"
)

// Service represents company application interface
type Service interface {
	Create(echo.Context, sandpiper.Company) (*sandpiper.Company, error)
	List(echo.Context, *sandpiper.Pagination) ([]sandpiper.Company, error)
	View(echo.Context, uuid.UUID) (*sandpiper.Company, error)
	Delete(echo.Context, uuid.UUID) error
	Update(echo.Context, *Update) (*sandpiper.Company, error)
}

// New creates new company application service
func New(db *pg.DB, sdb Repository, rbac RBAC, sec Securer) *Company {
	return &Company{db: db, sdb: sdb, rbac: rbac, sec: sec}
}

// Initialize initializes Company application service with defaults
func Initialize(db *pg.DB, rbac RBAC, sec Securer) *Company {
	return New(db, pgsql.NewCompany(), rbac, sec)
}

// Company represents company application service
type Company struct {
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
	Create(orm.DB, sandpiper.Company) (*sandpiper.Company, error)
	View(orm.DB, uuid.UUID) (*sandpiper.Company, error)
	List(orm.DB, *scope.Clause, *sandpiper.Pagination) ([]sandpiper.Company, error)
	Update(orm.DB, *sandpiper.Company) error
	Delete(orm.DB, *sandpiper.Company) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	CurrentUser(echo.Context) *sandpiper.AuthUser
	EnforceCompany(echo.Context, uuid.UUID) error
	EnforceRole(echo.Context, sandpiper.AccessLevel) error
}
