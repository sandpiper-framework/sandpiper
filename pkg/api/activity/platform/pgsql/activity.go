// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql

// activity service database access

import (
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"

	"autocare.org/sandpiper/pkg/shared/model"
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

	err := db.Model(activity).
		Column("activity.id", "message", "duration", "activity.created_at").
		Relation("Slice").WherePK().Select()
	if err != nil {
		return nil, err
	}
	return activity, nil
}

// CompanySubscribed checks if activity is included in a company's subscriptions.
func (s *Activity) CompanySubscribed(db orm.DB, companyID uuid.UUID, activityID uuid.UUID) bool {
	activity := new(sandpiper.Activity)
	err := db.Model(activity).Column("activity.id").
		Join("JOIN slices AS sl ON activity.slice_id = sl.id").
		Join("JOIN subscriptions AS sub ON sl.id = sub.slice_id").
		Where("sub.company_id = ?", companyID).
		Where("activity.id = ?", activityID).Select()
	if err == nil {
		return true
	}
	return false
}

// List returns a list of all activity with scoping and pagination
func (s *Activity) List(db orm.DB, sc *sandpiper.Scope, p *sandpiper.Pagination) ([]sandpiper.Activity, error) {
	var acts []sandpiper.Activity
	var err error

	// columns to select
	cols := "activity.id, message, duration, activity.created_at"

	if sc != nil {
		// Use CTE query to get all subscriptions for the scope (i.e. the company)
		err = db.Model((*sandpiper.Subscription)(nil)).
			Column("subscription.slice_id").
			Where(sc.Condition, sc.ID).
			WrapWith("scope").Table("scope").
			Join("JOIN activity ON activity.slice_id = scope.slice_id").
			ColumnExpr(cols).
			Limit(p.Limit).Offset(p.Offset).Select(&acts)
	} else {
		// simple case with no scoping
		err = db.Model(&acts).ColumnExpr(cols).Limit(p.Limit).Offset(p.Offset).Select()
	}
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
