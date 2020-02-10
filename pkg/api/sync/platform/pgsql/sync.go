// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql

// sync service database access

import (
	"autocare.org/sandpiper/pkg/shared/model"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
)

// Sync represents the client for sync table
type Sync struct{}

// NewSync returns a new sync database instance
func NewSync() *Sync {
	return &Sync{}
}

// Create creates a new sync in database (assumes allowed to do this)
func (s *Sync) Create(db orm.DB, sync sandpiper.Sync) (*sandpiper.Sync, error) {
	if err := db.Insert(&sync); err != nil {
		return nil, err
	}
	return &sync, nil
}

// View returns a single sync by ID (assumes allowed to do this)
func (s *Sync) View(db orm.DB, id int) (*sandpiper.Sync, error) {
	var sync = &sandpiper.Sync{ID: id}

	err := db.Model(sync).
		Column("sync.id", "message", "duration", "sync.created_at").
		Relation("Slice").WherePK().Select()
	if err != nil {
		return nil, err
	}
	return sync, nil
}

// CompanySubscribed checks if sync is included in a company's subscriptions.
func (s *Sync) CompanySubscribed(db orm.DB, companyID uuid.UUID, syncID uuid.UUID) bool {
	sync := new(sandpiper.Sync)
	err := db.Model(sync).Column("sync.id").
		Join("JOIN slices AS sl ON sync.slice_id = sl.id").
		Join("JOIN subscriptions AS sub ON sl.id = sub.slice_id").
		Where("sub.company_id = ?", companyID).
		Where("sync.id = ?", syncID).Select()
	if err == nil {
		return true
	}
	return false
}

// List returns a list of all syncs with scoping and pagination
func (s *Sync) List(db orm.DB, sc *sandpiper.Scope, p *sandpiper.Pagination) ([]sandpiper.Sync, error) {
	var syncs []sandpiper.Sync
	var err error

	// columns to select (optionally returning payload)
	cols := "sync.id, sync_type, sync_key, encoding, sync.created_at"

	if sc != nil {
		// Use CTE query to get all subscriptions for the scope (i.e. the company)
		err = db.Model((*sandpiper.Subscription)(nil)).
			Column("subscription.slice_id").
			Where(sc.Condition, sc.ID).
			WrapWith("scope").Table("scope").
			Join("JOIN syncs AS sync ON sync.slice_id = scope.slice_id").
			ColumnExpr(cols).
			Limit(p.Limit).Offset(p.Offset).Select(&syncs)
	} else {
		// simple case with no scoping
		err = db.Model(&syncs).ColumnExpr(cols).Limit(p.Limit).Offset(p.Offset).Select()
	}
	if err != nil {
		return nil, err
	}
	return syncs, nil
}

// Delete permanently removes a sync by primary key (id)
func (s *Sync) Delete(db orm.DB, id int) error {
	sync := sandpiper.Sync{ID: id}
	return db.Delete(sync)
}
