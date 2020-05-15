// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sandpiper

import (
	"github.com/google/uuid"
)

// AuthToken holds authentication token details with refresh token
type AuthToken struct {
	Token        string `json:"token"`
	Expires      string `json:"expires"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshToken holds authentication token details
type RefreshToken struct {
	Token   string `json:"token"`
	Expires string `json:"expires"`
}

// AuthUser represents data ("claims") stored in JWT token for the current user
type AuthUser struct {
	ID        int
	CompanyID uuid.UUID
	Username  string
	Email     string
	Role      AccessLevel
}

// APIKey is used to  authenticate a sync process
type APIKey struct {
	PrimaryID  uuid.UUID `json:"primary_id"`
	SyncAPIKey string    `json:"sync_api_key"`
}

// Server contains information about the current server
type Server struct {
	ID   uuid.UUID `json:"server-id"`
	Role string    `json:"server-role"`
}
