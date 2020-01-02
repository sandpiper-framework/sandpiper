// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sandpiper

import (
	"context"
	"time"

	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
)

// Subscription represents subscription model (m2m join table)
type Subscription struct {
	tableName    struct{}  `pg:"alias:sub"` // custom table alias for orm
	SliceID      uuid.UUID `json:"slice_id" pg:",pk"`
	SubscriberID uuid.UUID `json:"subscriber_id" pg:",pk"`
	Name         string    `json:"name"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at,omitempty" pg:",soft_delete"`
}

// compile-time check variables for model hooks (which take no memory)
var _ orm.BeforeInsertHook = (*Company)(nil)
var _ orm.BeforeUpdateHook = (*Company)(nil)

// BeforeInsert hooks into insert operations, setting createdAt and updatedAt to current time
func (b *Subscription) BeforeInsert(ctx context.Context) (context.Context, error) {
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
	return ctx, nil
}

// BeforeUpdate hooks into update operations, setting updatedAt to current time
func (b *Subscription) BeforeUpdate(ctx context.Context) (context.Context, error) {
	b.UpdatedAt = time.Now()
	return ctx, nil
}

func init() {
	// Register many to many model so ORM can better recognize m2m relation.
	// This should be done before dependant models are used.
	orm.RegisterTable((*Subscription)(nil))
}