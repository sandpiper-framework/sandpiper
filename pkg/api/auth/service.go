// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package auth

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"

	"sandpiper/pkg/api/auth/platform/pgsql"
	"sandpiper/pkg/shared/database"
	"sandpiper/pkg/shared/model"
	"sandpiper/pkg/shared/secure"
)

// Auth represents auth application service
type Auth struct {
	db   *pg.DB
	sdb  Repository
	tg   TokenGenerator
	sec  Securer
	rbac RBAC
}

// New creates a new auth service
func New(db *database.DB, sdb Repository, j TokenGenerator, sec Securer, rbac RBAC) *Auth {
	return &Auth{
		db:   db.DB,
		sdb:  sdb,
		tg:   j,
		sec:  sec,
		rbac: rbac,
	}
}

// Initialize initializes auth application service
func Initialize(db *database.DB, j TokenGenerator, sec Securer, rbac RBAC) *Auth {
	return New(db, pgsql.NewUser(), j, sec, rbac)
}

// Service represents auth service interface
type Service interface {
	Authenticate(echo.Context, string, string) (*sandpiper.AuthToken, error)
	Refresh(echo.Context, string) (*sandpiper.RefreshToken, error)
	Me(echo.Context) (*sandpiper.User, error)
	Server(echo.Context) *sandpiper.Server
	ParseCredentials(echo.Context) (*secure.Credentials, error)
}

// Repository represents available resource actions using a repository-abstraction-pattern interface.
type Repository interface {
	View(orm.DB, int) (*sandpiper.User, error)
	FindByUsername(orm.DB, string) (*sandpiper.User, error)
	FindByToken(orm.DB, string) (*sandpiper.User, error)
	Update(orm.DB, *sandpiper.User) error
}

// TokenGenerator represents token generator (jwt) interface
type TokenGenerator interface {
	GenerateToken(*sandpiper.User) (string, string, error)
}

// Securer represents security interface
type Securer interface {
	HashMatchesPassword(string, string) bool
	Token(string) string
	APIKeySecret() string
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	CurrentUser(echo.Context) *sandpiper.AuthUser
	OurServer() *sandpiper.Server
}
