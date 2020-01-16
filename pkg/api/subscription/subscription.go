// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package subscription contains services for the subscriptions resource.
package subscription

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/internal/model"
	"autocare.org/sandpiper/pkg/internal/scope"
)

// Create adds a new subscription if administrator
func (s *Subscription) Create(c echo.Context, req sandpiper.Subscription) (*sandpiper.Subscription, error) {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return nil, err
	}
	return s.sdb.Create(s.db, req)
}

// List returns list of subscriptions that you can view
func (s *Subscription) List(c echo.Context, p *sandpiper.Pagination) ([]sandpiper.Subscription, error) {
	au := s.rbac.CurrentUser(c)
	q, err := scope.Limit(au)
	if err != nil {
		return nil, err
	}
	return s.sdb.List(s.db, q, p)
}

// View returns a single subscription if allowed
func (s *Subscription) View(c echo.Context, sub sandpiper.Subscription) (*sandpiper.Subscription, error) {
	if err := s.rbac.EnforceSubscription(c, sub); err != nil {
		return nil, err
	}
	return s.sdb.View(s.db, sub)
}

// Delete deletes a subscription if administrator
func (s *Subscription) Delete(c echo.Context, id int) error {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return err
	}
	sub := sandpiper.Subscription{ID: id}
	subscription, err := s.sdb.View(s.db, sub)
	if err != nil {
		return err
	}
	return s.sdb.Delete(s.db, subscription)
}

// Update contains subscription field request used for updating
type Update struct {
	ID          int
	SliceID     uuid.UUID
	CompanyID   uuid.UUID
	Name        string
	Description string
	Active      bool
}

// Update updates subscription information
func (s *Subscription) Update(c echo.Context, r *Update) (*sandpiper.Subscription, error) {
	sub := sandpiper.Subscription{
		ID:          r.ID,
		SliceID:     r.SliceID,
		CompanyID:   r.CompanyID,
		Name:        r.Name,
		Description: r.Description,
		Active:      r.Active,
	}
	if err := s.rbac.EnforceSubscription(c, sub); err != nil {
		return nil, err
	}
	err := s.sdb.Update(s.db, &sub)
	if err != nil {
		return nil, err
	}
	return s.sdb.View(s.db, sub)
}
