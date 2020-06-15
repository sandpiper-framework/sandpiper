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

// convert url query params to sql clauses for filtering, ordering, paging and related model inclusion

/*
Query Strings:
	?sort=title,asc&sort=zipcode:desc,city,&filter=lname:Johnson,age:39&include=user
	?page=2&limit=20  # limit to define the number of items returned in the response
*/

// Params hold the url query parameters
type Params struct {
	RawQuery string
	Filter   []string
	Sort     []string
	Include  []string
	Paging   *sandpiper.Pagination
}

// Parse is a constructor for the query Params structure
func Parse(c echo.Context) (*Params, error) {

	p := &Params{
		RawQuery: c.QueryString(),
		Paging:   sandpiper.NewPagination(), // paging is never nil because of the way we set those values
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
		case "include":
			p.Include = v
		case "page":
			p.Paging.SetPage(v)
		case "limit":
			p.Paging.SetLimit(v)
		}
	}
	return p, nil
}

// AddSort includes zero or more order clauses to an existing query
// a missing direction implies ascending (asc)
// e.g. ?sort=title:asc,lname:desc
func (p *Params) AddSort(q *orm.Query) (added bool) {
	// can have zero or more sort instructions
	for _, sort := range p.Sort {
		// each sort can have one or more comma-separated sort fields
		for _, f := range strings.Split(sort, ",") {
			// an optional colon separates the order direction, but postgresql requires a space
			s := strings.ReplaceAll(f, ":", " ")
			q.Order(s)
			added = true
		}
	}
	return added
}

// AddFilter includes zero or more where clauses to an existing query
// Conditions are field:value and only support equivalence (=) (for now)
// All conditions within a filter (and across filters) are ANDed together (for now)
// All comparison values are strings (which works because of postgresql's "automatically coerced" literals
// https://dba.stackexchange.com/questions/238983)
// e.g. ?filter=lname:Johnson, age:39&filter=role:200
func (p *Params) AddFilter(q *orm.Query) {
	// can have zero or more filters
	for _, filter := range p.Filter {
		// each filter can have one or more comma-separated conditions
		for _, f := range strings.Split(filter, ",") {
			// we currently only support exact equals
			i := strings.Index(f, ":")
			if i != -1 {
				field := strings.TrimSpace(f[:i])
				value := strings.TrimSpace(f[i+1:])
				q.Where(field+" = ?", value)
			}
		}
	}
}

// WantRelated checks to see if a particular related model should be included in the results
func (p *Params) WantRelated(name string) bool {
	// can have zero or one include requests (stop after one match)
	for _, req := range p.Include {
		if strings.ToLower(req) == name {
			return true
		}
	}
	return false
}
