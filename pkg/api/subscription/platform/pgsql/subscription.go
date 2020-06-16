// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package pgsql

// subscription service database access

import (
	"net/http"
	"strings"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/params"
)

// Subscription represents the client for subscription table
type Subscription struct{}

// NewSubscription returns a new subscription database instance
func NewSubscription() *Subscription {
	return &Subscription{}
}

// Custom errors
var (
	ErrAlreadyExists      = echo.NewHTTPError(http.StatusInternalServerError, "Subscription name already exists.")
	ErrMissingQueryParams = echo.NewHTTPError(http.StatusInternalServerError, "no query params supplied")
)

// Create creates a new Subscription in database (assumes allowed to do this)
func (s *Subscription) Create(db orm.DB, sub sandpiper.Subscription) (*sandpiper.Subscription, error) {
	// don't add if duplicate name
	if err := checkDuplicate(db, sub.Name); err != nil {
		return nil, err
	}
	if err := db.Insert(&sub); err != nil {
		return nil, err
	}
	return &sub, nil
}

// View returns a single subscription by ID (assumes allowed to do this)
func (s *Subscription) View(db orm.DB, sub sandpiper.Subscription) (*sandpiper.Subscription, error) {
	var q *orm.Query

	// support several ways to query the subscription
	switch {
	case sub.SubID != uuid.Nil:
		q = queryByPrimaryKey(db, &sub)
	case sub.Name != "":
		q = queryByName(db, &sub)
	case sub.SliceID != uuid.Nil && sub.CompanyID != uuid.Nil:
		q = queryByJunctionKeys(db, &sub)
	default:
		return nil, ErrMissingQueryParams
	}

	err := q.Select()
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

// List returns list of all subscriptions
func (s *Subscription) List(db orm.DB, sc *sandpiper.Scope, p *params.Params) (subs []sandpiper.Subscription, err error) {

	q := queryAll(db, &subs).Limit(p.Paging.PageSize).Offset(p.Paging.Offset())
	if sc != nil {
		q.Where(sc.Condition, sc.ID)
	}
	p.AddFilter(q)
	p.AddSort(q, "name")
	p.Paging.Count, err = q.SelectAndCount()
	if err != nil {
		return nil, err
	}

	return subs, nil
}

// Update updates subscription info by primary key (assumes allowed to do this)
func (s *Subscription) Update(db orm.DB, sub *sandpiper.Subscription) error {
	_, err := db.Model(sub).UpdateNotZero()
	return err
}

// Delete removes the subscription by primary key
func (s *Subscription) Delete(db orm.DB, sub *sandpiper.Subscription) error {
	return db.Delete(sub)
}

// checkDuplicate returns true if name found in database
func checkDuplicate(db orm.DB, name string) error {
	// attempt to select by unique key
	m := new(sandpiper.Subscription)
	err := db.Model(m).
		Column("sub_id").
		Where("lower(name) = ?", strings.ToLower(name)).
		Select()

	switch err {
	case pg.ErrNoRows: // ok to add
		return nil
	case nil: // found a row, so a duplicate
		return ErrAlreadyExists
	default: // return any other problem found
		return err
	}
}

// queryAll returns a query for all subscriptions (including company and slice)
func queryAll(db orm.DB, subs *[]sandpiper.Subscription) *orm.Query {
	return db.Model(subs).Relation("Company").Relation("Slice")
}

// queryByPrimaryKey returns a query for a subscription by primary key (including company and slice)
func queryByPrimaryKey(db orm.DB, sub *sandpiper.Subscription) *orm.Query {
	return db.Model(sub).Relation("Company").Relation("Slice").WherePK()
}

// queryByJunctionKeys returns a query for a subscription by junction keys (including company and slice)
func queryByJunctionKeys(db orm.DB, sub *sandpiper.Subscription) *orm.Query {
	return db.Model(sub).Relation("Company").Relation("Slice").
		Where("slice_id = ? and company_id = ?", sub.SliceID, sub.CompanyID)
}

// queryByName returns a query for a subscription by name (including company and slice)
func queryByName(db orm.DB, sub *sandpiper.Subscription) *orm.Query {
	return db.Model(sub).Relation("Company").Relation("Slice").
		Where("lower(name) = ?", strings.ToLower(sub.Name))
}
