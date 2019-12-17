package password

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/password/platform/pgsql"
	"autocare.org/sandpiper/pkg/model"
)

// Service represents password application interface
type Service interface {
	Change(echo.Context, int, string, string) error
}

// Password represents password application service
type Password struct {
	db   *pg.DB
	udb  UserDB
	rbac RBAC
	sec  Securer
}

// New creates new password application service
func New(db *pg.DB, udb UserDB, rbac RBAC, sec Securer) *Password {
	return &Password{
		db:   db,
		udb:  udb,
		rbac: rbac,
		sec:  sec,
	}
}

// Initialize initializes password application service with defaults
func Initialize(db *pg.DB, rbac RBAC, sec Securer) *Password {
	return New(db, pgsql.NewUser(), rbac, sec)
}

// UserDB represents user repository interface
type UserDB interface {
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
