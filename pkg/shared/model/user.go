// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package sandpiper

import (
	"context"
	"time"

	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
)

// User represents user domain model
type User struct {
	ID              int         `json:"id"`
	FirstName       string      `json:"first_name"`
	LastName        string      `json:"last_name"`
	Username        string      `json:"username"`
	Password        string      `json:"-"`
	Email           string      `json:"email"`
	Phone           string      `json:"phone,omitempty"`
	Active          bool        `json:"active"`
	LastLogin       time.Time   `json:"last_login,omitempty"`
	PasswordChanged time.Time   `json:"password_changed,omitempty"`
	Token           string      `json:"-"`
	Role            AccessLevel `json:"role"`
	CompanyID       uuid.UUID   `json:"company_id"` // belongs-to company
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

// ChangePassword updates user's password related fields
func (u *User) ChangePassword(hash string) {
	u.Password = hash
	u.PasswordChanged = time.Now()
}

// UpdateLastLogin updates last login field
func (u *User) UpdateLastLogin(token string) {
	u.Token = token
	u.LastLogin = time.Now()
}

// compile-time check variables for model hooks (which take no memory)
var _ orm.BeforeInsertHook = (*User)(nil)
var _ orm.BeforeUpdateHook = (*User)(nil)

// BeforeInsert hooks into insert operations, setting createdAt and updatedAt to current time
func (u *User) BeforeInsert(ctx context.Context) (context.Context, error) {
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	return ctx, nil
}

// BeforeUpdate hooks into update operations, setting updatedAt to current time
func (u *User) BeforeUpdate(ctx context.Context) (context.Context, error) {
	u.UpdatedAt = time.Now()
	return ctx, nil
}

// UsersPaginated defines the list response
type UsersPaginated struct {
	Users  []User      `json:"data"`
	Paging *Pagination `json:"paging"`
}
