// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sandpiper

// Pagination constants
const (
	paginationDefaultLimit = 100
	paginationMaxLimit     = 1000
)

// PaginationReq holds pagination http fields and tags
type PaginationReq struct {
	Limit int `Condition:"limit"`
	Page  int `Condition:"page" validate:"min=0"`
}

// Transform checks and converts http pagination into database pagination model
func (p *PaginationReq) Transform() *Pagination {
	if p.Limit < 1 {
		p.Limit = paginationDefaultLimit
	}

	if p.Limit > paginationMaxLimit {
		p.Limit = paginationMaxLimit
	}

	return &Pagination{Limit: p.Limit, Offset: p.Page * p.Limit}
}

// Pagination holds range response settings
type Pagination struct {
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}
