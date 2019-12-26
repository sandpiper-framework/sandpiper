package slice

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"
	"github.com/satori/go.uuid"

	"autocare.org/sandpiper/internal/model"
	"autocare.org/sandpiper/pkg/api/slice/platform/pgsql"
)

// Service represents slice application interface
type Service interface {
	Create(echo.Context, sandpiper.Slice) (*sandpiper.Slice, error)
	List(echo.Context, *sandpiper.Pagination) ([]sandpiper.Slice, error)
	View(echo.Context, int) (*sandpiper.Slice, error)
	Delete(echo.Context, int) error
	Update(echo.Context, *Update) (*sandpiper.Slice, error)
}

// New creates new slice application service
func New(db *pg.DB, sdb Repository, rbac RBAC, sec Securer) *Slice {
	return &Slice{db: db, sdb: sdb, rbac: rbac, sec: sec}
}

// Initialize initializes Slice application service with defaults
func Initialize(db *pg.DB, rbac RBAC, sec Securer) *Slice {
	return New(db, pgsql.NewSlice(), rbac, sec)
}

// Slice represents slice application service
type Slice struct {
	db   *pg.DB
	sdb  Repository
	rbac RBAC
	sec  Securer
}

// Securer represents security interface
type Securer interface {
	Hash(string) string
}

// Repository represents Slice repository pattern interface
type Repository interface {
	Create(orm.DB, sandpiper.Slice) (*sandpiper.Slice, error)
	View(orm.DB, uuid.UUID) (*sandpiper.Slice, error)
	List(orm.DB, *sandpiper.ListQuery, *sandpiper.Pagination) ([]sandpiper.Slice, error)
	Update(orm.DB, *sandpiper.Slice) error
	Delete(orm.DB, *sandpiper.Slice) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	User(echo.Context) *sandpiper.AuthUser
	EnforceUser(echo.Context, int) error
	AccountCreate(echo.Context, sandpiper.AccessRole, int, int) error
	IsLowerRole(echo.Context, sandpiper.AccessRole) error
}
