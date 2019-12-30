// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sandpiper

import (
	"context"
	"github.com/go-pg/pg/v9/orm"
	"time"

	"github.com/google/uuid"
)

// User represents user domain model
type User struct {
	ID              int        `json:"id"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	Username        string     `json:"username"`
	Password        string     `json:"-"`
	Email           string     `json:"email"`
	Phone           string     `json:"phone,omitempty"`
	Active          bool       `json:"active"`
	LastLogin       time.Time  `json:"last_login,omitempty"`
	PasswordChanged time.Time  `json:"last_password_change,omitempty"`
	Token           string     `json:"-"`
	Role            AccessRole `json:"role,omitempty"`
	CompanyID       uuid.UUID  `json:"company_id"` // belongs-to company
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// AuthUser represents data stored in JWT token for the current user
type AuthUser struct {
	ID        int
	CompanyID uuid.UUID
	Username  string
	Email     string
	Role      AccessRole
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
