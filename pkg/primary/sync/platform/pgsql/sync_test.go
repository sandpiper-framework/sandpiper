// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"autocare.org/sandpiper/pkg/primary/sync/platform/pgsql"
	"autocare.org/sandpiper/pkg/shared/mock"
	"autocare.org/sandpiper/pkg/shared/model"
)

func TestCreate(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		req      sandpiper.Sync
		wantData *sandpiper.Sync
	}{
		{
			name:    "Fail on insert duplicate ID",
			wantErr: true,
			req: sandpiper.Sync{
				ID:       mock.TestUUID(1),
				SliceID:  mock.TestUUID(1),
				Type:     "aces-file",
				Key:      "AAP Premium Brakes",
				Encoding: "raw",
				Payload:  "payload data",
			},
		},
		{
			name:    "Fail on slice_id not found",
			wantErr: true,
			req: sandpiper.Sync{
				ID:       mock.TestUUID(1),
				SliceID:  mock.TestUUID(0),
				Type:     "aces-file",
				Key:      "AAP Premium Brakes",
				Encoding: "raw",
				Payload:  "payload data",
			},
		},
		{
			name: "Success",
			req: sandpiper.Sync{
				ID:           mock.TestUUID(2),
				Key:          "AAP Premium Brakes",
				ContentHash:  mock.TestHash(1),
				ContentCount: 1,
				LastUpdate:   mock.TestTime(2019),
			},
			wantData: &sandpiper.Sync{
				ID:           mock.TestUUID(2),
				Key:          "AAP Premium Brakes",
				ContentHash:  mock.TestHash(1),
				ContentCount: 1,
				LastUpdate:   mock.TestTime(2019),
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Sync{})

	if err := mock.InsertMultiple(db, &sandpiper.Sync{
		ID:  mock.TestUUID(1),
		Key: "Acme Brakes"}, &cases[1].req); err != nil {
		t.Error(err)
	}

	mdb := pgsql.NewSync()

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
		wantData *sandpiper.Sync
	}{
		{
			name:    "VIEW Sync does not exist",
			wantErr: true,
			id:      mock.TestUUID(2),
		},
		{
			name: "VIEW Success",
			id:   mock.TestUUID(1),
			wantData: &sandpiper.Sync{
				ID:           mock.TestUUID(1),
				Key:          "AAP Premium Brakes",
				ContentHash:  mock.TestHash(1),
				ContentCount: 1,
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Sync{})

	if err := mock.InsertMultiple(db, &sandpiper.Sync{
		ID:  mock.TestUUID(1),
		Key: "Acme Brakes"}, cases[1].wantData); err != nil {
		t.Error(err)
	}

	udb := pgsql.NewSync()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			user, err := udb.View(db, tt.id)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				if user == nil {
					t.Errorf("response was nil due to: %v", err)
				} else {
					tt.wantData.CreatedAt = user.CreatedAt
					assert.Equal(t, tt.wantData, user)
				}
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
		wantData []sandpiper.Sync
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
				Condition: "id = ?",
			},
			wantData: []sandpiper.Sync{
				{
					ID:           mock.TestUUID(1),
					Key:          "Brakes",
					ContentHash:  mock.TestHash(1),
					ContentCount: 1,
					LastUpdate:   mock.TestTime(2019),
				},
				{
					ID:           mock.TestUUID(1),
					Key:          "Brakes",
					ContentHash:  mock.TestHash(1),
					ContentCount: 1,
					LastUpdate:   mock.TestTime(2019),
				},
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Sync{})

	mdb := pgsql.NewSync()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			users, err := mdb.List(db, tt.qp, tt.pg)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				for i, v := range users {
					tt.wantData[i].CreatedAt = v.CreatedAt
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
		usr      *sandpiper.Sync
		wantData *sandpiper.Sync
	}{
		{
			name: "Success",
			usr: &sandpiper.Sync{
				ID: mock.TestUUID(1),
			},
			wantData: &sandpiper.Sync{
				ID:           mock.TestUUID(1),
				Key:          "Brakes",
				ContentHash:  mock.TestHash(1),
				ContentCount: 1,
				LastUpdate:   mock.TestTime(2019),
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Sync{})

	if err := mock.InsertMultiple(db, &sandpiper.Sync{
		ID:  mock.TestUUID(1),
		Key: "Acme Brakes"}, cases[0].wantData); err != nil {
		t.Error(err)
	}

	mdb := pgsql.NewSync()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			err := mdb.Delete(db, tt.usr)
			assert.Equal(t, tt.wantErr, err != nil)

		})
	}
}
