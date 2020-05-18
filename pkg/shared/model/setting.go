// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sandpiper

import "github.com/google/uuid"

// Setting represents the setting domain model
type Setting struct {
	ID         bool      `json:"id" pg:",pk"`
	ServerRole string    `json:"server_role"`
	ServerID   uuid.UUID `json:"server_id"`
}
