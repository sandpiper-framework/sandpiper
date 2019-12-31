// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql_test

import (
	"github.com/google/uuid"
	"testing"

	"github.com/stretchr/testify/assert"

	"autocare.org/sandpiper/internal/model"
	"autocare.org/sandpiper/pkg/api/slice/platform/pgsql"
	"autocare.org/sandpiper/test/mock"
)

func TestCreate(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		req      sandpiper.Slice
		wantData *sandpiper.Slice
	}{
		{
			name:    "Slice Name already exists",
			wantErr: true,
			req: sandpiper.Slice{
				ID:           mock.TestUUID(1),
				Name:         "AAP Premium Brakes",
				ContentHash:  mock.TestHash(1),
				ContentCount: 1,
				LastUpdate:   mock.TestTime(2019),
			},
		},
		{
			name:    "Fail on insert duplicate ID",
			wantErr: true,
			req: sandpiper.Slice{
				ID:           mock.TestUUID(1),
				Name:         "AAP Premium Brakes",
				ContentHash:  mock.TestHash(1),
				ContentCount: 1,
				LastUpdate:   mock.TestTime(2019),
			},
		},
		{
			name: "Success",
			req: sandpiper.Slice{
				ID:           mock.TestUUID(2),
				Name:         "AAP Premium Brakes",
				ContentHash:  mock.TestHash(1),
				ContentCount: 1,
				LastUpdate:   mock.TestTime(2019),
			},
			wantData: &sandpiper.Slice{
				ID:           mock.TestUUID(2),
				Name:         "AAP Premium Brakes",
				ContentHash:  mock.TestHash(1),
				ContentCount: 1,
				LastUpdate:   mock.TestTime(2019),
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Slice{})

	if err := mock.InsertMultiple(db, &sandpiper.Slice{
		ID:     mock.TestUUID(1),
		Name:   "Acme Brakes"}, &cases[1].req); err != nil {
		t.Error(err)
	}

	mdb := pgsql.NewSlice()

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
		wantData *sandpiper.Slice
	}{
		{
			name:    "VIEW Slice does not exist",
			wantErr: true,
			id:      mock.TestUUID(2),
		},
		{
			name: "VIEW Success",
			id:   mock.TestUUID(1),
			wantData: &sandpiper.Slice{
				ID:           mock.TestUUID(1),
				Name:         "AAP Premium Brakes",
				ContentHash:  mock.TestHash(1),
				ContentCount: 1,
				LastUpdate:   mock.TestTime(2019),
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Slice{})

	if err := mock.InsertMultiple(db, &sandpiper.Slice{
		ID:     mock.TestUUID(1),
		Name:   "Acme Brakes"}, cases[1].wantData); err != nil {
		t.Error(err)
	}

	udb := pgsql.NewSlice()

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
		data     *sandpiper.Slice
		wantData *sandpiper.Slice
	}{
		{
			name: "UPDATE Success",
			data: &sandpiper.Slice{
				ID:           mock.TestUUID(1),
				Name:         "Brakes",
				ContentHash:  mock.TestHash(1),
				ContentCount: 1,
				LastUpdate:   mock.TestTime(2019),
			},
			wantData: &sandpiper.Slice{
				ID:           mock.TestUUID(1),
				Name:         "AAP Premium Brakes",
				ContentHash:  mock.TestHash(1),
				ContentCount: 2,
				LastUpdate:   mock.TestTime(2019),
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Slice{})

	if err := mock.InsertMultiple(db, &sandpiper.Slice{
		ID:     mock.TestUUID(1),
		Name:   "Acme Brakes"}, cases[0].data); err != nil {
		t.Error(err)
	}

	mdb := pgsql.NewSlice()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := mdb.Update(db, tt.wantData)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				comp := &sandpiper.Slice{ID: tt.data.ID}
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
		wantData []sandpiper.Slice
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
				ID:    mock.TestUUID(1),
				Condition: "id = ?",
			},
			wantData: []sandpiper.Slice{
				{
					ID:           mock.TestUUID(1),
					Name:         "Brakes",
					ContentHash:  mock.TestHash(1),
					ContentCount: 1,
					LastUpdate:   mock.TestTime(2019),
				},
				{
					ID:           mock.TestUUID(1),
					Name:         "Brakes",
					ContentHash:  mock.TestHash(1),
					ContentCount: 1,
					LastUpdate:   mock.TestTime(2019),
				},
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Slice{})

	mdb := pgsql.NewSlice()

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
		usr      *sandpiper.Slice
		wantData *sandpiper.Slice
	}{
		{
			name: "Success",
			usr: &sandpiper.Slice{
				ID:        mock.TestUUID(1),
				DeletedAt: mock.TestTime(2018),
			},
			wantData: &sandpiper.Slice{
				ID:           mock.TestUUID(1),
				Name:         "Brakes",
				ContentHash:  mock.TestHash(1),
				ContentCount: 1,
				LastUpdate:   mock.TestTime(2019),
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Slice{})

	if err := mock.InsertMultiple(db, &sandpiper.Slice{
		ID:     mock.TestUUID(1),
		Name:   "Acme Brakes"}, cases[0].wantData); err != nil {
		t.Error(err)
	}

	mdb := pgsql.NewSlice()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			err := mdb.Delete(db, tt.usr)
			assert.Equal(t, tt.wantErr, err != nil)

			// Check if the deleted_at was set
		})
	}
}
