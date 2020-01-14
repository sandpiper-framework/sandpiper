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

// GrainType is an enum defining grain content
type GrainType int16

const ( // Never change these values, only add to list
	Grain_aces_file       GrainType = 1
	Grain_aces_item       GrainType = 2
	Grain_asset_file      GrainType = 3
	Grain_partspro_file   GrainType = 4
	Grain_partspro_item   GrainType = 5
	Grain_pies_file       GrainType = 6
	Grain_pies_item       GrainType = 7
	Grain_pies_marketcopy GrainType = 8
	Grain_pies_pricesheet GrainType = 9
)

// EncodingMethod is an enum describing how the payload is encoded
type EncodingMethod int16

const ( // Never change these values, only add to list
	EncRaw     EncodingMethod = 1
	EncB64     EncodingMethod = 2
	EncGzipB64 EncodingMethod = 3
)

// Grain represents the sandpiper syncable-object
type Grain struct {
	ID        uuid.UUID      `json:"id"`
	SliceID   uuid.UUID      `json:"slice_id"`
	Type      GrainType      `json:"grain_type"`
	Key       string         `json:"grain_key"`
	Payload   PayloadData    `json:"payload"`
	Encoding  EncodingMethod `json:"payload"`
	CreatedAt time.Time      `json:"created_at"`
	Slice     *Slice         `json:"slice"` // has-one relation
}

// compile-time check variables for model hooks (which take no memory)
var _ orm.BeforeInsertHook = (*Grain)(nil)

// BeforeInsert hooks into insert operations, setting createdAt and updatedAt to current time
func (b *Grain) BeforeInsert(ctx context.Context) (context.Context, error) {
	b.CreatedAt = time.Now()
	return ctx, nil
}
