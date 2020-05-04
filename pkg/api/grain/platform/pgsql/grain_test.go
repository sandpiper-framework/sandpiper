// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"sandpiper/pkg/api/grain/platform/pgsql"
	"sandpiper/pkg/shared/mock"
	"sandpiper/pkg/shared/model"
)

func TestCreate(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		req      sandpiper.Grain
		wantData *sandpiper.Grain
	}{
		{
			name:    "Fail on insert duplicate ID",
			wantErr: true,
			req: sandpiper.Grain{
				ID:       mock.TestUUID(1),
				SliceID:  mock.TestUUIDp(1),
				Key:      "AAP Premium Brakes",
				Encoding: "raw",
				Payload:  sandpiper.PayloadData("payload data"),
			},
		},
		{
			name:    "Fail on slice_id not found",
			wantErr: true,
			req: sandpiper.Grain{
				ID:       mock.TestUUID(1),
				SliceID:  mock.TestUUIDp(0),
				Key:      "AAP Premium Brakes",
				Encoding: "raw",
				Payload:  sandpiper.PayloadData("payload data"),
			},
		},
		{
			name: "Success",
			req: sandpiper.Grain{
				ID:        mock.TestUUID(2),
				Key:       "AAP Premium Brakes",
				CreatedAt: mock.TestTime(2019),
			},
			wantData: &sandpiper.Grain{
				ID:        mock.TestUUID(2),
				Key:       "AAP Premium Brakes",
				CreatedAt: mock.TestTime(2019),
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Grain{})

	if err := mock.InsertMultiple(db, &sandpiper.Grain{
		ID:  mock.TestUUID(1),
		Key: "Acme Brakes"}, &cases[1].req); err != nil {
		t.Error(err)
	}

	mdb := pgsql.NewGrain()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := mdb.Create(db, false, &tt.req)
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
		wantData *sandpiper.Grain
	}{
		{
			name:    "VIEW Grain does not exist",
			wantErr: true,
			id:      mock.TestUUID(2),
		},
		{
			name: "VIEW Success",
			id:   mock.TestUUID(1),
			wantData: &sandpiper.Grain{
				ID:  mock.TestUUID(1),
				Key: "AAP Premium Brakes",
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Grain{})

	if err := mock.InsertMultiple(db, &sandpiper.Grain{
		ID:  mock.TestUUID(1),
		Key: "Acme Brakes"}, cases[1].wantData); err != nil {
		t.Error(err)
	}

	udb := pgsql.NewGrain()

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
		wantData []sandpiper.Grain
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
			wantData: []sandpiper.Grain{
				{
					ID:        mock.TestUUID(1),
					Key:       "Brakes",
					CreatedAt: mock.TestTime(2019),
				},
				{
					ID:        mock.TestUUID(1),
					Key:       "Brakes",
					CreatedAt: mock.TestTime(2019),
				},
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Grain{})

	mdb := pgsql.NewGrain()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			grains, err := mdb.List(db, uuid.Nil, false, tt.qp, tt.pg)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				for i, v := range grains {
					tt.wantData[i].CreatedAt = v.CreatedAt
				}
				assert.Equal(t, tt.wantData, grains)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		grain    *sandpiper.Grain
		wantData *sandpiper.Grain
	}{
		{
			name: "Success",
			grain: &sandpiper.Grain{
				ID: mock.TestUUID(1),
			},
			wantData: &sandpiper.Grain{
				ID:        mock.TestUUID(1),
				Key:       "Brakes",
				CreatedAt: mock.TestTime(2019),
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Grain{})

	if err := mock.InsertMultiple(db, &sandpiper.Grain{
		ID:  mock.TestUUID(1),
		Key: "Acme Brakes"}, cases[0].wantData); err != nil {
		t.Error(err)
	}

	mdb := pgsql.NewGrain()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			err := mdb.Delete(db, tt.grain.ID)
			assert.Equal(t, tt.wantErr, err != nil)

		})
	}
}
