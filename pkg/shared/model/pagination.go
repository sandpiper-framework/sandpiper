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

	return &Pagination{Page: p.Page, Offset: p.Page * p.Limit, Limit: p.Limit}
}

// Pagination holds range response settings
type Pagination struct {
	Page   int `json:"page"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Count  int `json:"count"`
}
