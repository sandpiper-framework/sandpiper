// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package params

import (
	"net/url"
	"strings"

	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
)

// convert url query params to pg where clauses for filtering and ordering

/*
Query Strings:
	?sort=title,asc&filter=title:"bar"
	?sort=title,asc&sort=lname,desc&range={0,24}&filter={"title":"bar"}&expand=user
	filter=lname: "Johnson", age: 39
	page=2
  limit=20  # limit to define the number of items returned in the response
  sort = [title,asc lname,desc]
*/

// Params hold the url query parameters
type Params struct {
	RawQuery string
	Filter   []string
	Sort     []string
	Expand   []string
	Paging   *sandpiper.Pagination
}

// Parse is a constructor for the query Params structure
func Parse(c echo.Context) (*Params, error) {

	p := &Params{
		RawQuery: c.QueryString(),
		Paging:   &sandpiper.Pagination{}, // paging is never nil because of the way we set those values
	}

	params, err := url.ParseQuery(p.RawQuery) // c.QueryParams() does not return an err
	if err != nil {
		return nil, err
	}
	for k, v := range params {
		switch strings.ToLower(k) {
		case "filter":
			p.Filter = v
		case "sort":
			p.Sort = v
		case "expand":
			p.Expand = v
		case "page":
			p.Paging.SetPage(v)
		case "limit":
			p.Paging.SetLimit(v)
		}
	}
	return p, nil
}

// AddSort includes zero or more order clauses to an existing query
func (p *Params) AddSort(q *orm.Query) {
	for _, sort := range p.Sort {
		s := strings.ReplaceAll(sort, ",", " ")
		q.Order(s)
	}
}

// AddFilter includes zero or more and-ed where clauses to an existing query
// e.g. lname: "Johnson", age: 39
func (p *Params) AddFilter(q *orm.Query) {
	for _, filter := range p.Filter {
		for _, f := range strings.Split(filter, ",") {
			i := strings.Index(f, ":")
			if i != -1 {
				lhs := strings.TrimSpace(f[:i])
				rhs := strings.TrimSpace(f[i+1:])
				q.Where(lhs + " = ?", rhs)
			}
		}
	}
}
