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

//GrainType is an enum defining grain content
type GrainType int16

const ( // Never change this ordering!
	aces_file GrainType = iota
	aces_item
	asset_file
	partspro_file
	partspro_item
	pies_file
	pies_item
	pies_marketcopy
	pies_pricesheet
)

// Grain represents the sandpiper syncable-object
type Grain struct {
	ID        uuid.UUID   `json:"id"`
	SliceID   uuid.UUID   `json:"slice_id"`
	Type      GrainType   `json:"grain_type"`
	Key       string      `json:"grain_key"`
	Payload   PayloadData `json:"payload"`
	CreatedAt time.Time   `json:"created_at"`
	Slice     *Slice      `json:"slice"` // has-one relation
}

// compile-time check variables for model hooks (which take no memory)
var _ orm.BeforeInsertHook = (*Grain)(nil)

// BeforeInsert hooks into insert operations, setting createdAt and updatedAt to current time
func (b *Grain) BeforeInsert(ctx context.Context) (context.Context, error) {
	b.CreatedAt = time.Now()
	return ctx, nil
}
