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

// MetaMap is used for slice metadata serialization
type MetaMap map[string]string

// Slice represents a single slice container
type Slice struct {
	ID           uuid.UUID  `json:"id" pg:",pk"`
	Name         string     `json:"slice_name"`
	ContentHash  string     `json:"content_hash"`
	ContentCount uint       `json:"content_count"`
	ContentDate  time.Time  `json:"content_date"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Metadata     MetaMap    `json:"metadata,omitempty" pg:"-"`
	Companies    []*Company `json:"companies,omitempty" pg:"many2many:subscriptions"`
	//Grains       []*Grain  // has-many relation
}

// compile-time check variables for model hooks (which take no memory)
var _ orm.BeforeInsertHook = (*Slice)(nil)
var _ orm.BeforeUpdateHook = (*Slice)(nil)

// BeforeInsert hooks into insert operations, setting createdAt and updatedAt to current time
func (b *Slice) BeforeInsert(ctx context.Context) (context.Context, error) {
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
	return ctx, nil
}

// BeforeUpdate hooks into update operations, setting updatedAt to current time
func (b *Slice) BeforeUpdate(ctx context.Context) (context.Context, error) {
	b.UpdatedAt = time.Now()
	return ctx, nil
}

// SliceArray is an array of slices
type SliceArray []Slice

// The IDs method creates an array of slice_ids
func (a SliceArray) IDs() []uuid.UUID {
	var ids = make([]uuid.UUID, 0, len(a))

	for _, slice := range a {
		ids = append(ids, slice.ID)
	}
	return ids
}

// SliceMetadata contains information about a slice
type SliceMetadata struct {
	SliceID uuid.UUID `json:"-" pg:",pk"`
	Key     string    `json:"key" pg:",pk"`
	Value   string    `json:"val"`
}

// MetaArray is an array of slice metadata
type MetaArray []SliceMetadata

// The ToMap method converts array of metadata key/value structs to a map
func (a MetaArray) ToMap(sliceID uuid.UUID) MetaMap {
	mm := make(MetaMap)
	for _, meta := range a {
		if meta.SliceID == sliceID {
			mm[meta.Key] = meta.Value
		}
	}
	return mm
}
