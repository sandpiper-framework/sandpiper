// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sandpiper

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
)

// Tag allows grouping of slices
type Tag struct {
	ID          int       `json:"id" pg:",pk"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Slices      []*Slice  `json:"slices,omitempty" pg:"many2many:slice_tags"`
}

// SliceTag represents the many-to-many junction table
type SliceTag struct {
	TagID   int       `json:"-" pg:",pk"`
	SliceID uuid.UUID `json:"slice_id" pg:",pk"`
}

// CleanName removes invalid values from tag name (returning error if now empty)
func (t *Tag) CleanName() error {
	r := strings.NewReplacer(",", "_", " ", "_", "'", "", "\"", "")
	s := r.Replace(t.Name)
	if len(s) == 0 {
		return errors.New("invalid tag name")
	}
	t.Name = strings.ToLower(s)
	return nil
}

// compile-time check variables for model hooks (which take no memory)
var _ orm.BeforeInsertHook = (*Tag)(nil)
var _ orm.BeforeUpdateHook = (*Tag)(nil)

// BeforeInsert hooks into insert operations, setting createdAt and updatedAt to current time
func (b *Tag) BeforeInsert(ctx context.Context) (context.Context, error) {
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
	return ctx, nil
}

// BeforeUpdate hooks into update operations, setting updatedAt to current time
func (b *Tag) BeforeUpdate(ctx context.Context) (context.Context, error) {
	b.UpdatedAt = time.Now()
	return ctx, nil
}

func init() {
	// Register many to many model so ORM can better recognize m2m relation.
	// This should be done before dependant models are used.
	orm.RegisterTable((*SliceTag)(nil))
}
