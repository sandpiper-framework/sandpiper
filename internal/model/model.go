package sandpiper

import (
	"context"
	"time"

	"github.com/go-pg/pg/v9/orm"
)

// Base contains common fields for all tables
type Base struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty" pg:",soft_delete"`
}

// ListQuery holds company/location data used for list db queries
type ListQuery struct {
	Query string
	ID    int
}

// compile-time check variables for model hooks (which take no memory)
var _ orm.BeforeInsertHook = (*Base)(nil)
var _ orm.BeforeUpdateHook = (*Base)(nil)

// BeforeInsert hooks into insert operations, setting createdAt and updatedAt to current time
func (b *Base) BeforeInsert(ctx context.Context) (context.Context, error) {
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
	return ctx, nil
}

// BeforeUpdate hooks into update operations, setting updatedAt to current time
func (b *Base) BeforeUpdate(ctx context.Context) (context.Context, error) {
	b.UpdatedAt = time.Now()
	return ctx, nil
}