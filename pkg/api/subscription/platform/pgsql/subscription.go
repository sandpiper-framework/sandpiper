// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql

// subscription service database access

import (
	"net/http"
	"strings"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/shared/model"
)

// Subscription represents the client for subscription table
type Subscription struct{}

// NewSubscription returns a new subscription database instance
func NewSubscription() *Subscription {
	return &Subscription{}
}

// Custom errors
var (
	// ErrAlreadyExists indicates the subscription name is already used
	ErrAlreadyExists = echo.NewHTTPError(http.StatusInternalServerError, "Subscription name already exists.")
)

// Create creates a new Subscription in database (assumes allowed to do this)
func (s *Subscription) Create(db orm.DB, sub sandpiper.Subscription) (*sandpiper.Subscription, error) {
	// don't add if duplicate name
	if nameExists(db, sub.Name) {
		return nil, ErrAlreadyExists
	}
	if err := db.Insert(&sub); err != nil {
		return nil, err
	}
	return &sub, nil
}

// View returns a single subscription by ID (assumes allowed to do this)
func (s *Subscription) View(db orm.DB, sub sandpiper.Subscription) (*sandpiper.Subscription, error) {
	var q *orm.Query

	if sub.SubID != 0 {
		q = queryByPrimaryKey(db, &sub)
	} else {
		q = queryByJunctionKeys(db, &sub)
	}
	err := q.Select()
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

// List returns list of all subscriptions
func (s *Subscription) List(db orm.DB, sc *sandpiper.Scope, p *sandpiper.Pagination) ([]sandpiper.Subscription, error) {
	var subs []sandpiper.Subscription

	q := queryAll(db, &subs).Limit(p.Limit).Offset(p.Offset).Order("name")
	if sc != nil {
		q.Where(sc.Condition, sc.ID)
	}
	err := q.Select()
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

// nameExists returns true if name found in database
func nameExists(db orm.DB, name string) bool {
	m := new(sandpiper.Subscription)
	err := db.Model(m).Where("lower(name) = ?", strings.ToLower(name)).Select()
	return err != pg.ErrNoRows
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
