// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package subscription

// subscription service

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/api/subscription/platform/pgsql"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/database"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
)

// Service represents subscription application interface
type Service interface {
	Create(echo.Context, sandpiper.Subscription) (*sandpiper.Subscription, error)
	List(echo.Context, *sandpiper.Pagination) ([]sandpiper.Subscription, error)
	View(echo.Context, sandpiper.Subscription) (*sandpiper.Subscription, error)
	Delete(echo.Context, uuid.UUID) error
	Update(echo.Context, *Update) (*sandpiper.Subscription, error)
}

// New creates new company application service
func New(db *database.DB, sdb Repository, rbac RBAC, sec Securer) *Subscription {
	return &Subscription{db: db.DB, sdb: sdb, rbac: rbac, sec: sec}
}

// Initialize initializes Subscription application service with defaults
func Initialize(db *database.DB, rbac RBAC, sec Securer) *Subscription {
	return New(db, pgsql.NewSubscription(), rbac, sec)
}

// Subscription represents company application service
type Subscription struct {
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
	Create(orm.DB, sandpiper.Subscription) (*sandpiper.Subscription, error)
	View(orm.DB, sandpiper.Subscription) (*sandpiper.Subscription, error)
	List(orm.DB, *sandpiper.Scope, *sandpiper.Pagination) ([]sandpiper.Subscription, error)
	Update(orm.DB, *sandpiper.Subscription) error
	Delete(orm.DB, *sandpiper.Subscription) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	CurrentUser(echo.Context) *sandpiper.AuthUser
	EnforceCompany(echo.Context, uuid.UUID) error
	EnforceRole(echo.Context, sandpiper.AccessLevel) error
	EnforceScope(echo.Context) (*sandpiper.Scope, error)
}
