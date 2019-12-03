package pgsql

// user service database access

import (
	"net/http"
	"strings"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/model"
)

// NewUser returns a new user database instance
func NewUser() *User {
	return &User{}
}

// User represents the client for user table
type User struct{}

// Custom errors
var (
	ErrAlreadyExists = echo.NewHTTPError(http.StatusInternalServerError, "Username or email already exists.")
)

// Create creates a new user on database
func (u *User) Create(db orm.DB, usr sandpiper.User) (*sandpiper.User, error) {
	var user = new(sandpiper.User)
	err := db.Model(user).Where("lower(username) = ? or lower(email) = ? and deleted_at is null",
		strings.ToLower(usr.Username), strings.ToLower(usr.Email)).Select()

	if err != nil && err != pg.ErrNoRows {
		return nil, ErrAlreadyExists
	}

	if err := db.Insert(&usr); err != nil {
		return nil, err
	}
	return &usr, nil
}

// View returns single user by ID
func (u *User) View(db orm.DB, id int) (*sandpiper.User, error) {
	var user = new(sandpiper.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."id" = ? and deleted_at is null)`
	_, err := db.QueryOne(user, sql, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Update updates user's contact info
func (u *User) Update(db orm.DB, user *sandpiper.User) error {
	_, err := db.Model(user).UpdateNotZero()
	return err
}

// List returns list of all users retrievable for the current user, depending on role
func (u *User) List(db orm.DB, qp *sandpiper.ListQuery, p *sandpiper.Pagination) ([]sandpiper.User, error) {
	var users []sandpiper.User
	q := db.Model(&users).Column("user.*").Relation("Role").Limit(p.Limit).Offset(p.Offset).Where("deleted_at is null").Order("user.id desc")
	if qp != nil {
		q.Where(qp.Query, qp.ID)
	}
	if err := q.Select(); err != nil {
		return nil, err
	}
	return users, nil
}

// Delete sets deleted_at for a user
func (u *User) Delete(db orm.DB, user *sandpiper.User) error {
	return db.Delete(user)
}
