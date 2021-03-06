// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package sandpiper

import (
	"context"
	"strings"
	"time"

	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"

	"github.com/sandpiper-framework/sandpiper/pkg/shared/payload"
)

const (
	// L1GrainKey is always the same for level-1 grains
	L1GrainKey = "level-1"
)

// Grain represents the sandpiper syncable-object
type Grain struct {
	ID         uuid.UUID           `json:"id" pg:",pk"`
	SliceID    *uuid.UUID          `json:"slice_id,omitempty"` // must be pointer for omitempty to work here!
	Key        string              `json:"grain_key" pg:"grain_key"`
	Source     string              `json:"source"`
	Encoding   string              `json:"encoding"`
	PayloadLen int                 `json:"payload_len" pg:"-"` // calculated: "length(payload) AS payload_len"
	Payload    payload.PayloadData `json:"payload,omitempty"`
	CreatedAt  time.Time           `json:"created_at"`
	Slice      *Slice              `json:"slice,omitempty"` // has-one relation
}

// compile-time check variables for model hooks (which take no memory)
var _ orm.BeforeInsertHook = (*Grain)(nil)

// BeforeInsert hooks into insert operations, setting createdAt and updatedAt to current time
func (g *Grain) BeforeInsert(ctx context.Context) (context.Context, error) {
	g.CreatedAt = time.Now()
	return ctx, nil
}

// Display returns basic grain information as a string
func (g *Grain) Display() string {
	s := strings.Builder{}
	s.WriteString("grain_id: " + g.ID.String() + "\n")
	s.WriteString("slice_id: " + g.SliceID.String() + "\n")
	s.WriteString("key: \"" + g.Key + "\"\n")
	s.WriteString("source: \"" + g.Source + "\"\n")
	s.WriteString("created: " + g.CreatedAt.String() + "\n")
	return s.String()
}

// DisplayFull returns basic grain information plus decoded payload as a string
func (g *Grain) DisplayFull() string {
	var data string

	p, err := g.Payload.Decode(g.Encoding)
	if err != nil {
		data = "(" + err.Error() + ")"
	} else {
		data = p
	}
	s := strings.Builder{}
	s.WriteString("grain_id: " + g.ID.String() + "\n")
	s.WriteString("slice_id: " + g.SliceID.String() + "\n")
	s.WriteString("key: \"" + g.Key + "\"\n")
	s.WriteString("source: \"" + g.Source + "\"\n")
	s.WriteString("Payload: " + data + "\n")
	s.WriteString("created: " + g.CreatedAt.String() + "\n")
	return s.String()
}

// GrainsPaginated adds pagination
type GrainsPaginated struct {
	Grains []Grain     `json:"data"`
	Paging *Pagination `json:"paging"`
}
