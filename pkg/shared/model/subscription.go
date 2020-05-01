// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sandpiper

import (
	"context"
	"time"

	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
)

// Subscription represents subscription model (also a m2m junction table between companies and slices)
type Subscription struct {
	SubID       uuid.UUID `json:"id" pg:",pk"`
	SliceID     uuid.UUID `json:"-" pg:",unique:altkey"`
	CompanyID   uuid.UUID `json:"-" pg:",unique:altkey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Company     *Company  `json:"company"`
	Slice       *Slice    `json:"slice"`
}

// SubsPaginated adds pagination
type SubsPaginated struct {
	Subs []Subscription `json:"subs"`
	Page int            `json:"page"`
}

// compile-time check variables for model hooks (which take no memory)
var _ orm.BeforeInsertHook = (*Subscription)(nil)
var _ orm.BeforeUpdateHook = (*Subscription)(nil)

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
