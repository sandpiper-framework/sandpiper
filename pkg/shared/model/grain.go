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

// Grain represents the sandpiper syncable-object
type Grain struct {
	ID        uuid.UUID   `json:"id" pg:",pk"`
	SliceID   *uuid.UUID  `json:"slice_id,omitempty"` // must be pointer for omitempty to work here!
	Type      string      `json:"grain_type" pg:"grain_type"`
	Key       string      `json:"grain_key" pg:"grain_key"`
	Encoding  string      `json:"encoding"`
	Payload   PayloadData `json:"payload,omitempty"`
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
