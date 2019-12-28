// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"autocare.org/sandpiper/internal/model"
	"autocare.org/sandpiper/pkg/api/user/platform/pgsql"
	"autocare.org/sandpiper/test/mock"
)

func TestCreate(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		req      sandpiper.User
		wantData *sandpiper.User
	}{
		{
			name:    "User already exists",
			wantErr: true,
			req: sandpiper.User{
				Email:    "johndoe@mail.com",
				Username: "johndoe",
			},
		},
		{
			name:    "Fail on insert duplicate ID",
			wantErr: true,
			req: sandpiper.User{
				ID:        1,
				Email:     "tomjones@mail.com",
				FirstName: "Tom",
				LastName:  "Jones",
				Username:  "tomjones",
				Role:      sandpiper.SuperAdminRole,
				CompanyID: mock.TestUUID(1),
				Password:  "pass",
			},
		},
		{
			name: "Success",
			req: sandpiper.User{
				ID:        2,
				Email:     "newtomjones@mail.com",
				FirstName: "Tom",
				LastName:  "Jones",
				Username:  "newtomjones",
				Role:      sandpiper.SuperAdminRole,
				CompanyID: mock.TestUUID(1),
				Password:  "pass",
			},
			wantData: &sandpiper.User{
				ID:        2,
				Email:     "newtomjones@mail.com",
				FirstName: "Tom",
				LastName:  "Jones",
				Username:  "newtomjones",
				Role:      sandpiper.SuperAdminRole,
				CompanyID: mock.TestUUID(1),
				Password:  "pass",
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.User{})

	udb := pgsql.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := udb.Create(db, tt.req)
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
		id       int
		wantData *sandpiper.User
	}{
		{
			name:    "User does not exist",
			wantErr: true,
			id:      1000,
		},
		{
			name: "Success",
			id:   2,
			wantData: &sandpiper.User{
				ID:        2,
				Email:     "tomjones@mail.com",
				FirstName: "Tom",
				LastName:  "Jones",
				Username:  "tomjones",
				Role:      sandpiper.SuperAdminRole,
				CompanyID: mock.TestUUID(1),
				Password:  "newPass",
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.User{})

	udb := pgsql.NewUser()

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
		usr      *sandpiper.User
		wantData *sandpiper.User
	}{
		{
			name: "Success",
			usr: &sandpiper.User{
				ID:        2,
				FirstName: "Z",
				LastName:  "Freak",
				Phone:     "123456",
				Username:  "newUsername",
			},
			wantData: &sandpiper.User{
				ID:        2,
				Email:     "tomjones@mail.com",
				FirstName: "Z",
				LastName:  "Freak",
				Username:  "tomjones",
				Role:      sandpiper.SuperAdminRole,
				CompanyID: mock.TestUUID(1),
				Password:  "newPass",
				Phone:     "123456",
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.User{})

	udb := pgsql.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := udb.Update(db, tt.wantData)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				user := &sandpiper.User{
					ID: tt.usr.ID,
				}
				if err := db.Select(user); err != nil {
					t.Error(err)
				}
				tt.wantData.UpdatedAt = user.UpdatedAt
				tt.wantData.CreatedAt = user.CreatedAt
				tt.wantData.LastLogin = user.LastLogin
				assert.Equal(t, tt.wantData, user)
			}
		})
	}
}

func TestList(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		qp       *sandpiper.ListQuery
		pg       *sandpiper.Pagination
		wantData []sandpiper.User
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
			qp: &sandpiper.ListQuery{
				ID:    mock.TestUUID(1),
				Query: "company_id = ?",
			},
			wantData: []sandpiper.User{
				{
					ID:        2,
					Email:     "tomjones@mail.com",
					FirstName: "Tom",
					LastName:  "Jones",
					Username:  "tomjones",
					Role:      sandpiper.SuperAdminRole,
					CompanyID: mock.TestUUID(1),
					Password:  "newPass",
				},
				{
					ID:        1,
					Email:     "johndoe@mail.com",
					FirstName: "John",
					LastName:  "Doe",
					Username:  "johndoe",
					Role:      sandpiper.SuperAdminRole,
					CompanyID: mock.TestUUID(1),
					Password:  "hunter2",
					Token:     "loginrefresh",
				},
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Role{}, &sandpiper.User{})

	if err := mock.InsertMultiple(db, &sandpiper.Role{
		ID:          1,
		AccessLevel: 1,
		Name:        "SUPER_ADMIN"}, &cases[1].wantData); err != nil {
		t.Error(err)
	}

	udb := pgsql.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			users, err := udb.List(db, tt.qp, tt.pg)
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
		usr      *sandpiper.User
		wantData *sandpiper.User
	}{
		{
			name: "Success",
			usr: &sandpiper.User{
				ID: 2,
			},
			wantData: &sandpiper.User{
				ID:        2,
				Email:     "tomjones@mail.com",
				FirstName: "Tom",
				LastName:  "Jones",
				Username:  "tomjones",
				Role:      sandpiper.SuperAdminRole,
				CompanyID: mock.TestUUID(1),
				Password:  "newPass",
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.Role{}, &sandpiper.User{})

	if err := mock.InsertMultiple(db, &sandpiper.Role{
		ID:          1,
		AccessLevel: 1,
		Name:        "SUPER_ADMIN"}, cases[0].wantData); err != nil {
		t.Error(err)
	}

	udb := pgsql.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			err := udb.Delete(db, tt.usr)
			assert.Equal(t, tt.wantErr, err != nil)

			// Check if the deleted_at was set
		})
	}
}
