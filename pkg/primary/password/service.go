// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package password

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/primary/password/platform/pgsql"
	"autocare.org/sandpiper/pkg/shared/model"
)

// Service represents password application interface
type Service interface {
	Change(echo.Context, int, string, string) error
}

// Password represents password application service
type Password struct {
	db   *pg.DB
	sdb  Repository
	rbac RBAC
	sec  Securer
}

// New creates new password application service
func New(db *pg.DB, sdb Repository, rbac RBAC, sec Securer) *Password {
	return &Password{
		db:   db,
		sdb:  sdb,
		rbac: rbac,
		sec:  sec,
	}
}

// Initialize initializes password application service with defaults
func Initialize(db *pg.DB, rbac RBAC, sec Securer) *Password {
	return New(db, pgsql.NewUser(), rbac, sec)
}

// Repository represents available resource actions using a repository-abstraction-pattern interface.
type Repository interface {
	View(orm.DB, int) (*sandpiper.User, error)
	Update(orm.DB, *sandpiper.User) error
}

// Securer represents security interface
type Securer interface {
	Hash(string) string
	HashMatchesPassword(string, string) bool
	Password(string, ...string) bool
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	EnforceUser(echo.Context, int) error
}