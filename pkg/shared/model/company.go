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

// Company represents company model
type Company struct {
	ID            uuid.UUID       `json:"id"`
	Name          string          `json:"name"`
	SyncAddr      string          `json:"sync_addr"`
	SyncAPIKey    string          `json:"sync_api_key,omitempty"` // only on secondary
	SyncUserID    int             `json:"sync_user_id,omitempty"` // only on primary
	Active        bool            `json:"active"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	Users         []*User         `json:"users,omitempty"`         // has-many relation
	SyncUser      *User           `json:"sync_user,omitempty"`     // has-one relation
	Subscriptions []*Subscription `json:"subscriptions,omitempty"` // has-many relation
}

// compile-time check variables for model hooks (which take no memory)
var _ orm.BeforeInsertHook = (*Company)(nil)
var _ orm.BeforeUpdateHook = (*Company)(nil)

// BeforeInsert hooks into insert operations, setting createdAt and updatedAt to current time
func (b *Company) BeforeInsert(ctx context.Context) (context.Context, error) {
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
	return ctx, nil
}

// BeforeUpdate hooks into update operations, setting updatedAt to current time
func (b *Company) BeforeUpdate(ctx context.Context) (context.Context, error) {
	b.UpdatedAt = time.Now()
	return ctx, nil
}

// CompaniesPaginated adds pagination
type CompaniesPaginated struct {
	Companies []Company   `json:"data"`
	Paging    *Pagination `json:"paging"`
}
