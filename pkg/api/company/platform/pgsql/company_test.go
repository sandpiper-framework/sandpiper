// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql_test

import (
	"github.com/google/uuid"
	"testing"

	"github.com/stretchr/testify/assert"

	"autocare.org/sandpiper/pkg/api/company/platform/pgsql"
	"autocare.org/sandpiper/pkg/internal/model"
	"autocare.org/sandpiper/test/mock"
)

func TestCreate(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		req      sandpiper.Company
		wantData *sandpiper.Company
	}{
		{
			name:    "CREATE Company Name already exists",
			wantErr: true,
			req: sandpiper.Company{
				ID:     mock.TestUUID(1),
				Name:   "Acme Brakes",
				Active: true,
			},
		},
		{
			name:    "CREATE Fail on insert duplicate ID",
			wantErr: true,
			req: sandpiper.Company{
				ID:     mock.TestUUID(1),
				Name:   "Acme Brakes",
				Active: true,
			},
		},
		{
			name: "CREATE Success",
			req: sandpiper.Company{
				ID:     mock.TestUUID(2),
				Name:   "Acme Brakes",
				Active: true,
			},
			wantData: &sandpiper.Company{
				ID:     mock.TestUUID(2),
				Name:   "Acme Brakes",
				Active: true,
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Company{})

	if err := mock.InsertMultiple(db, &sandpiper.Company{
		ID:     mock.TestUUID(1),
		Name:   "Acme Brakes",
		Active: true}, &cases[1].req); err != nil {
		t.Error(err)
	}

	mdb := pgsql.NewCompany()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := mdb.Create(db, tt.req)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				if resp == nil {
					t.Error("Expected data, but received nil.")
					return
				}
				tt.wantData.CreatedAt = resp.CreatedAt
				tt.wantData.UpdatedAt = resp.UpdatedAt
				assert.Equal(t, tt.wantData, resp)
			}
		})
	}
}

func TestView(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		id       uuid.UUID
		wantData *sandpiper.Company
	}{
		{
			name:    "VIEW Company does not exist",
			wantErr: true,
			id:      mock.TestUUID(2),
		},
		{
			name: "VIEW Success",
			id:   mock.TestUUID(1),
			wantData: &sandpiper.Company{
				ID:     mock.TestUUID(1),
				Name:   "Acme Brakes",
				Active: true,
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Company{})

	if err := mock.InsertMultiple(db, &sandpiper.Company{
		ID:     mock.TestUUID(1),
		Name:   "Acme Brakes",
		Active: true}, cases[1].wantData); err != nil {
		t.Error(err)
	}

	udb := pgsql.NewCompany()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			user, err := udb.View(db, tt.id)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				if user == nil {
					t.Errorf("response was nil due to: %v", err)
				} else {
					tt.wantData.CreatedAt = user.CreatedAt
					tt.wantData.UpdatedAt = user.UpdatedAt
					assert.Equal(t, tt.wantData, user)
				}
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		data     *sandpiper.Company
		wantData *sandpiper.Company
	}{
		{
			name: "UPDATE Success",
			data: &sandpiper.Company{
				ID:     mock.TestUUID(1),
				Name:   "Before Update",
				Active: false,
			},
			wantData: &sandpiper.Company{
				ID:     mock.TestUUID(1),
				Name:   "Acme Brakes",
				Active: true,
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Company{})

	if err := mock.InsertMultiple(db, &sandpiper.Company{
		ID:     mock.TestUUID(1),
		Name:   "Acme Brakes",
		Active: true}, cases[0].data); err != nil {
		t.Error(err)
	}

	mdb := pgsql.NewCompany()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := mdb.Update(db, tt.wantData)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				comp := &sandpiper.Company{ID: tt.data.ID}
				if err := db.Select(comp); err != nil {
					t.Error(err)
				}
				tt.wantData.UpdatedAt = comp.UpdatedAt
				tt.wantData.CreatedAt = comp.CreatedAt
				tt.wantData.DeletedAt = comp.DeletedAt
				assert.Equal(t, tt.wantData, comp)
			}
		})
	}
}

func TestList(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		qp       *scope.Clause
		pg       *sandpiper.Pagination
		wantData []sandpiper.Company
	}{
		{
			name:    "Invalid pagination values",
			wantErr: true,
			pg: &sandpiper.Pagination{
				Limit: -100,
			},
		},
		{
			name: "Success",
			pg: &sandpiper.Pagination{
				Limit:  100,
				Offset: 0,
			},
			qp: &scope.Clause{
				ID:        mock.TestUUID(1),
				Condition: "company_id = ?",
			},
			wantData: []sandpiper.Company{
				{
					ID:     mock.TestUUID(1),
					Name:   "Acme Brakes",
					Active: true,
				},
				{
					ID:     mock.TestUUID(2),
					Name:   "Acme Wipers",
					Active: true,
				},
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Company{})

	mdb := pgsql.NewCompany()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			users, err := mdb.List(db, tt.qp, tt.pg)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				for i, v := range users {
					tt.wantData[i].CreatedAt = v.CreatedAt
					tt.wantData[i].UpdatedAt = v.UpdatedAt
				}
				assert.Equal(t, tt.wantData, users)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		usr      *sandpiper.Company
		wantData *sandpiper.Company
	}{
		{
			name: "DELETE Success",
			usr: &sandpiper.Company{
				ID:        mock.TestUUID(1),
				DeletedAt: mock.TestTime(2018),
			},
			wantData: &sandpiper.Company{
				ID:     mock.TestUUID(1),
				Name:   "Acme Brakes",
				Active: true,
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Company{})

	if err := mock.InsertMultiple(db, &sandpiper.Company{
		ID:     mock.TestUUID(1),
		Name:   "Acme Brakes",
		Active: true}, cases[0].wantData); err != nil {
		t.Error(err)
	}

	mdb := pgsql.NewCompany()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			err := mdb.Delete(db, tt.usr)
			assert.Equal(t, tt.wantErr, err != nil)

			// Check if the deleted_at was set
		})
	}
}
