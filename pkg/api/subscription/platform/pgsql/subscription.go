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

	"autocare.org/sandpiper/pkg/internal/model"
	"autocare.org/sandpiper/pkg/internal/scope"
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
	if sub.ID != 0 {
		return findByPrimaryKey(db, sub)
	}
	return findByJunction(db, sub)
}

// Update updates subscription info by primary key (assumes allowed to do this)
func (s *Subscription) Update(db orm.DB, sub *sandpiper.Subscription) error {
	_, err := db.Model(sub).Update()
	return err
}

// List returns list of all subscriptions
func (s *Subscription) List(db orm.DB, sc *scope.Clause, p *sandpiper.Pagination) ([]sandpiper.Subscription, error) {
	var subs []sandpiper.Subscription

	q := db.Model(&subs).Limit(p.Limit).Offset(p.Offset).Order("name")
	if sc != nil {
		q.Where(sc.Condition, sc.ID)
	}
	if err := q.Select(); err != nil {
		return nil, err
	}
	return subs, nil
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

// findByPrimaryKey returns a subscription using supplied primary key
func findByPrimaryKey(db orm.DB, sub sandpiper.Subscription) (*sandpiper.Subscription, error) {
	err := db.Select(&sub)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

// findByJunction returns a subscription using supplied junction table keys
func findByJunction(db orm.DB, sub sandpiper.Subscription) (*sandpiper.Subscription, error) {
	err := db.Model(&sub).Where("slice_id = ? and company_id = ?", sub.SliceID, sub.CompanyID).Select()
	if err != nil {
		return nil, err
	}
	return &sub, nil
}
