// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sandpiper

import (
	"github.com/google/uuid"
)

// ListQuery holds company data used for list db queries
type ListQuery struct {
	Query string
	ID    uuid.UUID
}

