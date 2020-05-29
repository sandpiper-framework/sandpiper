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

// Activity logs for sync requests
type Activity struct {
	tableName    struct{}      `pg:"activity"` // we don't want the plural `activities`
	ID           int           `json:"id" pg:",pk"`
	SubID        uuid.UUID     `json:"sub_id"`
	Success      bool          `json:"success" pg:",use_zero"`
	Message      string        `json:"message"`
	Error        string        `json:"error"`
	Duration     time.Duration `json:"duration"`
	CreatedAt    time.Time     `json:"created_at"`
	Subscription *Subscription `json:"subscription"` // has one
}

// ActivityPaginated adds pagination
type ActivityPaginated struct {
	Syncs []Activity `json:"activity"`
	Page  int        `json:"page"`
}

// compile-time check variables for model hooks (which take no memory)
var _ orm.BeforeInsertHook = (*Activity)(nil)

// BeforeInsert hooks into insert operations, setting createdAt and updatedAt to current time
func (b *Activity) BeforeInsert(ctx context.Context) (context.Context, error) {
	b.CreatedAt = time.Now()
	return ctx, nil
}
