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
)

const (
	// PrimaryServer is a constant ServerRole value
	PrimaryServer = "primary"

	// SecondaryServer is a constant ServerRole value
	SecondaryServer = "secondary"
)

// Setting represents the setting domain model
type Setting struct {
	ID         bool      `json:"id" pg:",pk"`
	ServerRole string    `json:"server_role"`
	ServerID   uuid.UUID `json:"server_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// compile-time check variables for model hooks (which take no memory)
var _ orm.BeforeInsertHook = (*Setting)(nil)
var _ orm.BeforeUpdateHook = (*Setting)(nil)

// BeforeInsert hooks into insert operations, setting createdAt and updatedAt to current time
func (s *Setting) BeforeInsert(ctx context.Context) (context.Context, error) {
	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	return ctx, nil
}

// BeforeUpdate hooks into update operations, setting updatedAt to current time
func (s *Setting) BeforeUpdate(ctx context.Context) (context.Context, error) {
	s.UpdatedAt = time.Now()
	return ctx, nil
}

// Display returns basic grain information as a string
func (s *Setting) Display() string {
	b := strings.Builder{}
	b.WriteString("Server ID: " + s.ServerID.String() + "\n")
	b.WriteString("Server Role: " + s.ServerRole + "\n")
	b.WriteString("Created: " + s.CreatedAt.String() + "\n")
	b.WriteString("Updated: " + s.UpdatedAt.String() + "\n")
	return b.String()
}
