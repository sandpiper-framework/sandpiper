// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package auth

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/internal/model"
	"autocare.org/sandpiper/pkg/api/auth/platform/pgsql"
)

// Auth represents auth application service
type Auth struct {
	db   *pg.DB
	udb  Repository
	tg   TokenGenerator
	sec  Securer
	rbac RBAC
}

// New creates a new auth service
func New(db *pg.DB, udb Repository, j TokenGenerator, sec Securer, rbac RBAC) *Auth {
	return &Auth{
		db:   db,
		udb:  udb,
		tg:   j,
		sec:  sec,
		rbac: rbac,
	}
}

// Initialize initializes auth application service
func Initialize(db *pg.DB, j TokenGenerator, sec Securer, rbac RBAC) *Auth {
	return New(db, pgsql.NewUser(), j, sec, rbac)
}

// Service represents auth service interface
type Service interface {
	Authenticate(echo.Context, string, string) (*sandpiper.AuthToken, error)
	Refresh(echo.Context, string) (*sandpiper.RefreshToken, error)
	Me(echo.Context) (*sandpiper.User, error)
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
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	CurrentUser(echo.Context) *sandpiper.AuthUser
}
