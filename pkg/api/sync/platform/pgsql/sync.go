// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql

// sync service database access

import (
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"

	"autocare.org/sandpiper/pkg/shared/model"
)

// Sync represents the client for sync table
type Sync struct{}

// NewSync returns a new sync database instance
func NewSync() *Sync {
	return &Sync{}
}

// LogActivity permanently removes a sync by primary key (id)
func (s *Sync) LogActivity(db orm.DB, req sandpiper.SyncRequest) error {

	return nil
}

// Primary returns a single (primary) company by ID (assumes allowed to do this)
func (s *Sync) Primary(db orm.DB, id uuid.UUID) (*sandpiper.Company, error) {
	var company = &sandpiper.Company{ID: id}

	err := db.Model(company).WherePK().Select()
	if err != nil {
		return nil, err
	}
	return company, nil
}
