// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sandpiper

import (
	"net/url"
	"strings"
)

// TagQuery is used to store tag query parameters
type TagQuery struct {
	RawQuery string
	IsUnion  bool
	TagList  []string
}

// NewTagQuery returns a new tag query structure with parsed tags
func NewTagQuery(params url.Values, q string) *TagQuery {
	tq := new(TagQuery)
	tq.RawQuery = q
	for k, v := range params {
		tags := strings.ReplaceAll(v[0], " ", "")
		if k == "tags" {
			tq.IsUnion = true
			tq.TagList = strings.Split(tags, ",")
			return tq
		}
		if k == "tags-all" {
			tq.TagList = strings.Split(tags, ",")
			return tq
		}
	}
	return nil
}

// Provided checks to see if a tag query was included in the url
func (q *TagQuery) Provided() bool {
	return len(q.TagList) > 0
}

// Count returns the number of tags provided in the query
func (q *TagQuery) Count() int {
	return len(q.TagList)
}