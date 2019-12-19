package pgsql

// auth service database access

import (
	"github.com/go-pg/pg/v9/orm"

	"autocare.org/sandpiper/internal/model"
)

// User represents the client for user table
type User struct{}

// NewUser returns a new user database instance (for the service)
func NewUser() *User {
	return &User{}
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

// FindByUsername queries for single user by username
func (u *User) FindByUsername(db orm.DB, uname string) (*sandpiper.User, error) {
	var user = new(sandpiper.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."username" = ? and deleted_at is null)`
	_, err := db.QueryOne(user, sql, uname)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindByToken queries for single user by token
func (u *User) FindByToken(db orm.DB, token string) (*sandpiper.User, error) {
	var user = new(sandpiper.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."token" = ? and deleted_at is null)`
	_, err := db.QueryOne(user, sql, token)
	if err != nil {
	}
	return user, err
}

// Update updates user's info
func (u *User) Update(db orm.DB, user *sandpiper.User) error {
	return db.Update(user)
}
