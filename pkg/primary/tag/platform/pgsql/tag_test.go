// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"autocare.org/sandpiper/pkg/primary/tag/platform/pgsql"
	"autocare.org/sandpiper/pkg/shared/mock"
	"autocare.org/sandpiper/pkg/shared/model"
)

func TestCreate(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		req      sandpiper.Tag
		wantData *sandpiper.Tag
	}{
		{
			name:    "Tag Name already exists",
			wantErr: true,
			req: sandpiper.Tag{
				ID:     mock.TestUUID(1),
				Name:   "Acme Brakes",
				Active: true,
			},
		},
		{
			name:    "Fail on insert duplicate ID",
			wantErr: true,
			req: sandpiper.Tag{
				ID:     mock.TestUUID(1),
				Name:   "Acme Brakes",
				Active: true,
			},
		},
		{
			name: "Success",
			req: sandpiper.Tag{
				ID:     mock.TestUUID(2),
				Name:   "Acme Brakes",
				Active: true,
			},
			wantData: &sandpiper.Tag{
				ID:     mock.TestUUID(2),
				Name:   "Acme Brakes",
				Active: true,
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Tag{})

	if err := mock.InsertMultiple(db, &sandpiper.Tag{
		ID:     mock.TestUUID(1),
		Name:   "Acme Brakes",
		Active: true}, &cases[1].req); err != nil {
		t.Error(err)
	}

	mdb := pgsql.NewTag()

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
		wantData *sandpiper.Tag
	}{
		{
			name:    "Tag does not exist",
			wantErr: true,
			id:      mock.TestUUID(2),
		},
		{
			name: "Success",
			id:   mock.TestUUID(1),
			wantData: &sandpiper.Tag{
				ID:     mock.TestUUID(1),
				Name:   "Acme Brakes",
				Active: true,
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Tag{})

	if err := mock.InsertMultiple(db, &sandpiper.Tag{
		ID:     mock.TestUUID(1),
		Name:   "Acme Brakes",
		Active: true}, cases[1].wantData); err != nil {
		t.Error(err)
	}

	udb := pgsql.NewTag()

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
		data     *sandpiper.Tag
		wantData *sandpiper.Tag
	}{
		{
			name: "Success",
			data: &sandpiper.Tag{
				ID:     mock.TestUUID(1),
				Name:   "Before Update",
				Active: false,
			},
			wantData: &sandpiper.Tag{
				ID:     mock.TestUUID(1),
				Name:   "Acme Brakes",
				Active: true,
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Tag{})

	if err := mock.InsertMultiple(db, &sandpiper.Tag{
		ID:     mock.TestUUID(1),
		Name:   "Acme Brakes",
		Active: true}, cases[0].data); err != nil {
		t.Error(err)
	}

	mdb := pgsql.NewTag()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := mdb.Update(db, tt.wantData)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				comp := &sandpiper.Tag{ID: tt.data.ID}
				if err := db.Select(comp); err != nil {
					t.Error(err)
				}
				tt.wantData.UpdatedAt = comp.UpdatedAt
				tt.wantData.CreatedAt = comp.CreatedAt
				assert.Equal(t, tt.wantData, comp)
			}
		})
	}
}

func TestList(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		qp       *sandpiper.Scope
		pg       *sandpiper.Pagination
		wantData []sandpiper.Tag
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
			qp: &sandpiper.Scope{
				ID:        mock.TestUUID(1),
				Condition: "tag_id = ?",
			},
			wantData: []sandpiper.Tag{
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

	db := mock.NewDB(t, dbCon, &sandpiper.Tag{})

	mdb := pgsql.NewTag()

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
		usr      *sandpiper.Tag
		wantData *sandpiper.Tag
	}{
		{
			name: "Success",
			usr: &sandpiper.Tag{
				ID: mock.TestUUID(1),
			},
			wantData: &sandpiper.Tag{
				ID:     mock.TestUUID(1),
				Name:   "Acme Brakes",
				Active: true,
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Tag{})

	if err := mock.InsertMultiple(db, &sandpiper.Tag{
		ID:     mock.TestUUID(1),
		Name:   "Acme Brakes",
		Active: true}, cases[0].wantData); err != nil {
		t.Error(err)
	}

	mdb := pgsql.NewTag()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			err := mdb.Delete(db, tt.usr)
			assert.Equal(t, tt.wantErr, err != nil)

		})
	}
}