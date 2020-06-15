// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package params

import (
	"net/url"
	"strings"
)

// TagQuery is used to store tag query parameters
type TagQuery struct {
	IsUnion bool
	TagList []string
}

// NewTagQuery searches the url query string for tag filters
func NewTagQuery(tags url.Values) *TagQuery {
	tq := new(TagQuery)
	for k, v := range tags {
		if k == "tags" {
			tq.assign(v[0], true)
		}
		if k == "tags-all" {
			tq.assign(v[0], false)
		}
	}
	return tq
}

// assign adds values to the underlying struct
func (q *TagQuery) assign(vals string, isUnion bool) {
	q.TagList = strings.Split(vals, ",")
	q.IsUnion = isUnion
}

// Provided checks to see if a tag query was included in the url
func (q *TagQuery) Provided() bool {
	return len(q.TagList) > 0
}

// Count returns the number of tags provided in the query
func (q *TagQuery) Count() int {
	return len(q.TagList)
}
