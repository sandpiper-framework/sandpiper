package sandpiper

import (
	"time"

	"github.com/satori/go.uuid"
)

// Slice represents a single slice container
type Slice struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"slice_name"`
	ContentHash  string    `json:"content_hash"`
	ContentCount uint      `json:"content_count"`
	LastUpdate   time.Time `json:"last_update"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at,omitempty" pg:",soft_delete"`
}
