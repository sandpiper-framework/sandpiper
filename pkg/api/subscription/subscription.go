// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package subscription contains services for the subscriptions resource.
package subscription

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/params"
)

// Create adds a new subscription if administrator
func (s *Subscription) Create(c echo.Context, req sandpiper.Subscription) (*sandpiper.Subscription, error) {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return nil, err
	}
	return s.sdb.Create(s.db, req)
}

// List returns list of subscriptions that you can view
func (s *Subscription) List(c echo.Context, p *params.Params) ([]sandpiper.Subscription, error) {
	q, err := s.rbac.EnforceScope(c)
	if err != nil {
		return nil, err
	}
	return s.sdb.List(s.db, q, p)
}

// View returns a single subscription if allowed
func (s *Subscription) View(c echo.Context, sub sandpiper.Subscription) (*sandpiper.Subscription, error) {
	// must get it first to see if we are allowed to return it
	subscription, err := s.sdb.View(s.db, sub)
	if err != nil {
		return nil, err
	}
	if err := s.rbac.EnforceCompany(c, subscription.CompanyID); err != nil {
		return nil, err
	}
	return subscription, nil
}

// Delete deletes a subscription if administrator
func (s *Subscription) Delete(c echo.Context, id uuid.UUID) error {
	if err := s.rbac.EnforceRole(c, sandpiper.AdminRole); err != nil {
		return err
	}
	sub := sandpiper.Subscription{SubID: id}
	subscription, err := s.sdb.View(s.db, sub)
	if err != nil {
		return err
	}
	return s.sdb.Delete(s.db, subscription)
}

// Update contains subscription field request used for updating
type Update struct {
	SubID       uuid.UUID
	SliceID     uuid.UUID
	CompanyID   uuid.UUID
	Name        string
	Description string
	Active      bool
}

// Update updates subscription information
func (s *Subscription) Update(c echo.Context, r *Update) (*sandpiper.Subscription, error) {
	if err := s.rbac.EnforceCompany(c, r.CompanyID); err != nil {
		return nil, err
	}
	sub := sandpiper.Subscription{
		SubID:       r.SubID,
		SliceID:     r.SliceID,
		CompanyID:   r.CompanyID,
		Name:        r.Name,
		Description: r.Description,
		Active:      r.Active,
	}
	err := s.sdb.Update(s.db, &sub)
	if err != nil {
		return nil, err
	}
	return s.sdb.View(s.db, sub)
}
