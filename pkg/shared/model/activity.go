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

// Activity logs for sync requests
type Activity struct {
	tableName    struct{}      `pg:"activity"` // we don't want the plural `activities`
	ID           int           `json:"id" pg:",pk"`
	CompanyID    uuid.UUID     `json:"company_id"`
	SubID        uuid.UUID     `json:"sub_id"`
	Success      bool          `json:"success" pg:",use_zero"`
	Message      string        `json:"message"`
	Error        string        `json:"error"`
	Duration     time.Duration `json:"duration"`
	CreatedAt    time.Time     `json:"created_at"`
	Subscription *Subscription `json:"subscription,omitempty"` // has one
}

// ActivityPaginated adds pagination
type ActivityPaginated struct {
	Syncs  []Activity  `json:"activity"`
	Paging *Pagination `json:"paging"`
}

// compile-time check variables for model hooks (which take no memory)
var _ orm.BeforeInsertHook = (*Activity)(nil)

// BeforeInsert hooks into insert operations, setting createdAt and updatedAt to current time
func (b *Activity) BeforeInsert(ctx context.Context) (context.Context, error) {
	b.CreatedAt = time.Now()
	return ctx, nil
}
