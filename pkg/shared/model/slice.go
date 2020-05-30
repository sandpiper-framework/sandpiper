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

// MetaMap is used for slice metadata serialization
type MetaMap map[string]string

// SyncStatus Enum
const (
	SyncStatusNone     = "none"
	SyncStatusUpdating = "updating"
	SyncStatusSuccess  = "success"
	SyncStatusError    = "error"
)

// Slice represents a single slice container
type Slice struct {
	ID              uuid.UUID  `json:"id" pg:",pk"`
	Name            string     `json:"slice_name"`
	SliceType       string     `json:"slice_type"`
	ContentHash     string     `json:"content_hash"`
	ContentCount    int        `json:"content_count"`
	ContentDate     time.Time  `json:"content_date"`
	AllowSync       bool       `json:"allow_sync"`
	SyncStatus      string     `json:"sync_status"`
	LastSyncAttempt time.Time  `json:"last_sync_attempt"`
	LastGoodSync    time.Time  `json:"last_good_sync"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Metadata        MetaMap    `json:"metadata,omitempty" pg:"-"`
	Companies       []*Company `json:"companies,omitempty" pg:"many2many:subscriptions"`
}

// SlicesPaginated adds pagination
type SlicesPaginated struct {
	Slices []Slice `json:"slices"`
	Page   int     `json:"page"`
}

// compile-time check variables for model hooks (which take no memory)
var _ orm.BeforeInsertHook = (*Slice)(nil)
var _ orm.BeforeUpdateHook = (*Slice)(nil)

// BeforeInsert hooks into insert operations, setting createdAt and updatedAt to current time
func (s *Slice) BeforeInsert(ctx context.Context) (context.Context, error) {
	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	return ctx, nil
}

// BeforeUpdate hooks into update operations, setting updatedAt to current time
func (s *Slice) BeforeUpdate(ctx context.Context) (context.Context, error) {
	s.UpdatedAt = time.Now()
	return ctx, nil
}

// Validate checks slice_type enum
func (s Slice) Validate() bool {
	// using `select enum_range(null::slice_type_enum)`
	switch s.SliceType {
	case "aces-file",
		"aces-items",
		"asset-files",
		"pies-file",
		"pies-items",
		"pies-marketcopy",
		"pies-pricesheet",
		"partspro-file",
		"partspro-items":
		return true
	}
	return false
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

// Equals checks if two MetaMaps have identical key/value pairs
func (a MetaMap) Equals(b MetaMap) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v1 := range a {
		if v2, ok := b[k]; !ok || v1 != v2 {
			return false
		}
	}
	return true
}
