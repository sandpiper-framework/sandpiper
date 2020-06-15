// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package sandpiper

import (
	"strconv"
)

// Pagination constants
const (
	paginationDefaultLimit = 100
	paginationMaxLimit     = 1000
)

// Pagination holds range response settings
type Pagination struct {
	Page  int `json:"page_number"`
	Limit int `json:"page_size"`
	Count int `json:"items_total"`
}

// SetPage is a setter method for page number
func (p *Pagination) SetPage(pages []string) {
	if len(pages) > 0 {
		p.Page, _ = strconv.Atoi(pages[0])
	}
	if p.Page == 0 {
		p.Page = 1
	}
}

// SetLimit is a setter method for item limit per page
func (p *Pagination) SetLimit(limits []string) {
	if len(limits) > 0 {
		p.Limit, _ = strconv.Atoi(limits[0])
	}
	if p.Limit < 1 {
		p.Limit = paginationDefaultLimit
	}
	if p.Limit > paginationMaxLimit {
		p.Limit = paginationMaxLimit
	}
}

// Offset method calculates the sql offset value
func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}

// PaginationReq holds pagination http fields and tags
type PaginationReq struct {
	Limit int `query:"limit"`
	Page  int `query:"page" validate:"min=0"`
}
