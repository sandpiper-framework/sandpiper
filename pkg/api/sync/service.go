// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sync

import (
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"sandpiper/pkg/api/sync/platform/pgsql"
	"sandpiper/pkg/shared/client"
	"sandpiper/pkg/shared/database"
	"sandpiper/pkg/shared/model"
)

// Service represents sync application interface
type Service interface {
	Start(echo.Context, uuid.UUID) error
	Process(echo.Context) error
	Subscriptions(c echo.Context) ([]sandpiper.Subscription, error)
	Grains(echo.Context, uuid.UUID, bool) ([]sandpiper.Grain, error)
}

// New creates new sync application service
func New(db *database.DB, sdb Repository, rbac RBAC, sec Securer) *Sync {
	return &Sync{db: db.DB, sdb: sdb, rbac: rbac, sec: sec}
}

// Initialize initializes Sync application service with defaults
func Initialize(db *database.DB, rbac RBAC, sec Securer) *Sync {
	return New(db, pgsql.NewSync(), rbac, sec)
}

// Sync represents sync application service
type Sync struct {
	db   *pg.DB
	sdb  Repository
	rbac RBAC
	sec  Securer
	key  string         // secret key for en/decrypting sync credentials
	api  *client.Client // maintains our login to partner server for sync
}

// Securer represents security interface
type Securer interface {
	Hash(string) string
}

// Repository represents available resource actions using a repository-abstraction-pattern interface.
type Repository interface {
	Primary(orm.DB, uuid.UUID) (*sandpiper.Company, error)
	LogActivity(orm.DB, uuid.UUID, uuid.UUID, string, time.Duration, error) error
	Subscriptions(orm.DB, uuid.UUID) ([]sandpiper.Subscription, error)
	AddSubscription(orm.DB, sandpiper.Subscription) error
	DeactivateSubscription(orm.DB, uuid.UUID) error
	SliceAccess(orm.DB, uuid.UUID, uuid.UUID) error
	AddSlice(orm.DB, *sandpiper.Slice) error
	RefreshSlice(orm.DB, *sandpiper.Slice) error
	SliceMetadata(orm.DB, uuid.UUID) (sandpiper.MetaArray, error)
	ReplaceSliceMetadata(orm.DB, uuid.UUID, sandpiper.MetaArray) error
	Grains(orm.DB, uuid.UUID, bool) ([]sandpiper.Grain, error)
	AddGrain(orm.DB, *sandpiper.Grain) error
	DeleteGrains(orm.DB, []uuid.UUID) error
	BeginSyncUpdate(orm.DB, uuid.UUID) error
	FinalizeSyncUpdate(orm.DB, uuid.UUID, error) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	CurrentUser(echo.Context) *sandpiper.AuthUser
	EnforceRole(echo.Context, sandpiper.AccessLevel) error
	EnforceScope(echo.Context) (*sandpiper.Scope, error)
	EnforceServerRole(string) error
}
