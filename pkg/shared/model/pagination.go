// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

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

	if p.Page == 0 {
		p.Page = 1
	}

	return &Pagination{Page: p.Page, Limit: p.Limit, Offset: (p.Page - 1) * p.Limit}
}

// Pagination holds range response settings
type Pagination struct {
	Page   int `json:"page_number"`
	Limit  int `json:"items_limit"`
	Count  int `json:"items_total"`
	Offset int `json:"-"` // for database queries
}
