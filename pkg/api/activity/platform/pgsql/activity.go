// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package pgsql

// activity service database access

import (
	"github.com/go-pg/pg/v9/orm"

	"sandpiper/pkg/shared/model"
)

// Activity represents the client for activity table
type Activity struct{}

// NewActivity returns a new activity database instance
func NewActivity() *Activity {
	return &Activity{}
}

// Create creates a new activity in database (assumes allowed to do this)
func (s *Activity) Create(db orm.DB, activity sandpiper.Activity) (*sandpiper.Activity, error) {
	if err := db.Insert(&activity); err != nil {
		return nil, err
	}
	return &activity, nil
}

// View returns a single activity by ID (assumes allowed to do this)
func (s *Activity) View(db orm.DB, id int) (*sandpiper.Activity, error) {
	var activity = &sandpiper.Activity{ID: id}

	err := db.Model(activity).Relation("Subscription").WherePK().Select()
	if err != nil {
		return nil, err
	}
	return activity, nil
}

// List returns a list of all activity with scoping and pagination
func (s *Activity) List(db orm.DB, p *sandpiper.Pagination) ([]sandpiper.Activity, error) {
	var acts []sandpiper.Activity
	var err error

	err = db.Model(&acts).Relation("Subscription").Limit(p.Limit).Offset(p.Offset).Select()
	if err != nil {
		return nil, err
	}
	return acts, nil
}

// Delete permanently removes an activity by primary key (id)
func (s *Activity) Delete(db orm.DB, id int) error {
	activity := sandpiper.Activity{ID: id}
	return db.Delete(activity)
}
