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

// Subscription represents subscription model (also a m2m junction table between companies and slices)
type Subscription struct {
	SubID       uuid.UUID `json:"id" pg:",pk"`
	SliceID     uuid.UUID `json:"slice_id" pg:",unique:altkey"`
	CompanyID   uuid.UUID `json:"company_id" pg:",unique:altkey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Company     *Company  `json:"company,omitempty"`
	Slice       *Slice    `json:"slice,omitempty"`
}

// SubsPaginated adds pagination
type SubsPaginated struct {
	Subs   []Subscription `json:"subs"`
	Paging *Pagination    `json:"paging"`
}

// SubsMap allows fast lookups by SubID
type SubsMap map[uuid.UUID]Subscription

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

// SemiDeepCopy makes a copy of subscription without sharing memory (but one level deep)
func (b *Subscription) SemiDeepCopy() Subscription {
	sub := *b
	if sub.Company != nil {
		company := *sub.Company
		sub.Company = &company
	}
	if sub.Slice != nil {
		slice := *sub.Slice
		sub.Slice = &slice
	}
	return sub
}

func init() {
	// Register many to many model so ORM can better recognize m2m relation.
	// This should be done before dependant models are used.
	orm.RegisterTable((*Subscription)(nil))
}

// Load fills a map of subs for fast lookup [sub_id: sub]
func (m SubsMap) Load(subs []Subscription) {
	for _, sub := range subs {
		m[sub.SubID] = sub
	}
}
