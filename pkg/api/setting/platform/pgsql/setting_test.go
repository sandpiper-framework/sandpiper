// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package pgsql_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"sandpiper/pkg/api/setting/platform/pgsql"
	"sandpiper/pkg/shared/mock"
	"sandpiper/pkg/shared/model"
)

func TestCreate(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		req      sandpiper.Setting
		wantData *sandpiper.Setting
	}{
		{
			name:    "Fail on insert duplicate ID",
			wantErr: true,
			req: sandpiper.Setting{
				ID:       mock.TestUUID(1),
				SliceID:  mock.TestUUIDp(1),
				Key:      "AAP Premium Brakes",
				Encoding: "raw",
				Payload:  payload.PayloadData("payload data"),
			},
		},
		{
			name:    "Fail on slice_id not found",
			wantErr: true,
			req: sandpiper.Setting{
				ID:       mock.TestUUID(1),
				SliceID:  mock.TestUUIDp(0),
				Key:      "AAP Premium Brakes",
				Encoding: "raw",
				Payload:  payload.PayloadData("payload data"),
			},
		},
		{
			name: "Success",
			req: sandpiper.Setting{
				ID:        mock.TestUUID(2),
				Key:       "AAP Premium Brakes",
				CreatedAt: mock.TestTime(2019),
			},
			wantData: &sandpiper.Setting{
				ID:        mock.TestUUID(2),
				Key:       "AAP Premium Brakes",
				CreatedAt: mock.TestTime(2019),
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Setting{})

	if err := mock.InsertMultiple(db, &sandpiper.Setting{
		ID:  mock.TestUUID(1),
		Key: "Acme Brakes"}, &cases[1].req); err != nil {
		t.Error(err)
	}

	mdb := pgsql.NewSetting()

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
		wantData *sandpiper.Setting
	}{
		{
			name:    "VIEW Setting does not exist",
			wantErr: true,
			id:      mock.TestUUID(2),
		},
		{
			name: "VIEW Success",
			id:   mock.TestUUID(1),
			wantData: &sandpiper.Setting{
				ID:  mock.TestUUID(1),
				Key: "AAP Premium Brakes",
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Setting{})

	if err := mock.InsertMultiple(db, &sandpiper.Setting{
		ID:  mock.TestUUID(1),
		Key: "Acme Brakes"}, cases[1].wantData); err != nil {
		t.Error(err)
	}

	udb := pgsql.NewSetting()

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
		wantData []sandpiper.Setting
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
			wantData: []sandpiper.Setting{
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

	db := mock.NewDB(t, dbCon, &sandpiper.Setting{})

	mdb := pgsql.NewSetting()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			settings, err := mdb.List(db, uuid.Nil, false, tt.qp, tt.pg)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				for i, v := range settings {
					tt.wantData[i].CreatedAt = v.CreatedAt
				}
				assert.Equal(t, tt.wantData, settings)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		setting  *sandpiper.Setting
		wantData *sandpiper.Setting
	}{
		{
			name: "Success",
			setting: &sandpiper.Setting{
				ID: mock.TestUUID(1),
			},
			wantData: &sandpiper.Setting{
				ID:        mock.TestUUID(1),
				Key:       "Brakes",
				CreatedAt: mock.TestTime(2019),
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Setting{})

	if err := mock.InsertMultiple(db, &sandpiper.Setting{
		ID:  mock.TestUUID(1),
		Key: "Acme Brakes"}, cases[0].wantData); err != nil {
		t.Error(err)
	}

	mdb := pgsql.NewSetting()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			err := mdb.Delete(db, tt.setting.ID)
			assert.Equal(t, tt.wantErr, err != nil)

		})
	}
}
