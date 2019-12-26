package user

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/internal/model"
	"autocare.org/sandpiper/pkg/api/user/platform/pgsql"
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
	udb  UDB
	rbac RBAC
	sec  Securer
}

// New creates new user application service
func New(db *pg.DB, udb UDB, rbac RBAC, sec Securer) *User {
	return &User{db: db, udb: udb, rbac: rbac, sec: sec}
}

// Initialize initializes User application service with defaults
func Initialize(db *pg.DB, rbac RBAC, sec Securer) *User {
	return New(db, pgsql.NewUser(), rbac, sec)
}

// Securer represents security interface
type Securer interface {
	Hash(string) string
}

// UDB represents user repository interface
type UDB interface {
	Create(orm.DB, sandpiper.User) (*sandpiper.User, error)
	View(orm.DB, int) (*sandpiper.User, error)
	List(orm.DB, *sandpiper.ListQuery, *sandpiper.Pagination) ([]sandpiper.User, error)
	Update(orm.DB, *sandpiper.User) error
	Delete(orm.DB, *sandpiper.User) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	CurrentUser(echo.Context) *sandpiper.AuthUser
	EnforceUser(echo.Context, int) error
	AccountCreate(echo.Context, sandpiper.AccessRole, int, int) error
	IsLowerRole(echo.Context, sandpiper.AccessRole) error
}
