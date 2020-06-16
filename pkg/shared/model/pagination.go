// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package sandpiper

import (
	"strconv"
)

// limit the number of items returned from a list query and allow paging through them

// Pagination constants
const (
	defaultPageSize = 100
	maxPageSize     = 1000
)

// Pagination holds range response settings
type Pagination struct {
	PageNumber int `json:"page_number"`
	PageSize   int `json:"page_size"`
	Count      int `json:"items_total"`
}

// NewPagination is the Pagination constructor with default values
func NewPagination() *Pagination {
	return &Pagination{PageNumber: 1, PageSize: defaultPageSize}
}

// SetPageNumber is a setter method for page number
func (p *Pagination) SetPageNumber(pages []string) {
	if len(pages) > 0 {
		p.PageNumber, _ = strconv.Atoi(pages[0])
	}
	if p.PageNumber == 0 {
		p.PageNumber = 1
	}
}

// SetPageSize is a setter method for item limit per page
func (p *Pagination) SetPageSize(limits []string) {
	if len(limits) > 0 {
		p.PageSize, _ = strconv.Atoi(limits[0])
	}
	if p.PageSize < 1 {
		p.PageSize = defaultPageSize
	}
	if p.PageSize > maxPageSize {
		p.PageSize = maxPageSize
	}
}

// Offset method calculates the sql offset value
func (p *Pagination) Offset() int {
	return (p.PageNumber - 1) * p.PageSize
}
