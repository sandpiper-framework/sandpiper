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
	IsUnion bool
	TagList []string
}

// NewTagQuery returns a new tag query structure
func NewTagQuery(params url.Values) *TagQuery {
	tq := new(TagQuery)
	for k, v := range params {
		tags := strings.ReplaceAll(v[0]," ", "")
		if k == "tags" {
			tq.IsUnion = true
			tq.TagList = strings.Split(tags, ",")
			return tq
		}
		if k == "tags-any" {
			tq.TagList = strings.Split(tags, ",")
			return tq
		}
	}
	return nil
}
