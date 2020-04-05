// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sandpiper

import (
	"context"
	"strings"
	"time"

	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
)

// Grain represents the sandpiper syncable-object
type Grain struct {
	ID        uuid.UUID   `json:"id" pg:",pk"`
	SliceID   *uuid.UUID  `json:"slice_id,omitempty"` // must be pointer for omitempty to work here!
	Key       string      `json:"grain_key" pg:"grain_key"`
	Source    string      `json:"source"`
	Encoding  string      `json:"encoding"`
	Payload   PayloadData `json:"payload,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	Slice     *Slice      `json:"slice,omitempty"` // has-one relation
}

// compile-time check variables for model hooks (which take no memory)
var _ orm.BeforeInsertHook = (*Grain)(nil)

// BeforeInsert hooks into insert operations, setting createdAt and updatedAt to current time
func (b *Grain) BeforeInsert(ctx context.Context) (context.Context, error) {
	b.CreatedAt = time.Now()
	return ctx, nil
}

// Display returns basic grain information in string format
func (g *Grain) Display() string {
	s := strings.Builder{}
	s.WriteString("grain_id: " + g.ID.String() + "\n")
	s.WriteString("slice_id: " + g.SliceID.String() + "\n")
	s.WriteString("key: \"" + g.Key + "\"\n")
	s.WriteString("source: \"" + g.Source + "\"\n")
	s.WriteString("created: " + g.CreatedAt.String() + "\n")
	return s.String()
}

// DisplayFull returns abbreviated grain information in string format
func (g *Grain) DisplayFull() string {
	var payload string

	p, err := g.Payload.Decode()
	if err != nil {
		payload = "(" + err.Error() + ")"
	} else {
		payload = string(p)
	}
	s := strings.Builder{}
	s.WriteString("grain_id: " + g.ID.String() + "\n")
	s.WriteString("slice_id: " + g.SliceID.String() + "\n")
	s.WriteString("key: \"" + g.Key + "\"\n")
	s.WriteString("source: \"" + g.Source + "\"\n")
	s.WriteString("Payload: " + payload + "\n")
	s.WriteString("created: " + g.CreatedAt.String() + "\n")
	return s.String()
}

// GrainsPaginated adds pagination
type GrainsPaginated struct {
	Grains []Grain `json:"grains"`
	Page   int     `json:"page"`
}
