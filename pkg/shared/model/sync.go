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

// Sync logs sync requests
type Sync struct {
	ID        int        `json:"id" pg:",pk"`
	SliceID   *uuid.UUID `json:"slice_id,omitempty"` // must be pointer for omitempty to work here!
	Message   string     `json:"message"`
	Duration  time.Time  `json:"duration"`
	CreatedAt time.Time  `json:"created_at"`
	Slice     *Slice     `json:"slice,omitempty"` // has-one relation
}

// compile-time check variables for model hooks (which take no memory)
var _ orm.BeforeInsertHook = (*Sync)(nil)

// BeforeInsert hooks into insert operations, setting createdAt and updatedAt to current time
func (b *Sync) BeforeInsert(ctx context.Context) (context.Context, error) {
	b.CreatedAt = time.Now()
	return ctx, nil
}
