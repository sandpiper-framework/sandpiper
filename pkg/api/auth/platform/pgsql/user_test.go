// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"autocare.org/sandpiper/pkg/api/auth/platform/pgsql"
	"autocare.org/sandpiper/pkg/shared/mock"
	"autocare.org/sandpiper/pkg/shared/model"
)

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
			name:    "Success",
			wantErr: false,
			id:      1,
			wantData: &sandpiper.User{
				ID:        1,
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

	// seed with what we want to find
	if err := mock.InsertMultiple(db, cases[1].wantData); err != nil {
		t.Error(err)
	}

	udb := pgsql.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			user, err := udb.View(db, tt.id)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				if user == nil {
					t.Errorf("response was nil due to: %v", err)
				} else {
					tt.wantData.CreatedAt = user.CreatedAt // should be set by the insert
					tt.wantData.UpdatedAt = user.UpdatedAt // should be set by the insert
					assert.Equal(t, tt.wantData, user)
				}
			}
		})
	}
}

func TestFindByUsername(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		username string
		wantData *sandpiper.User
	}{
		{
			name:     "User does not exist",
			wantErr:  true,
			username: "notExists",
		},
		{
			name:     "Success",
			wantErr:  false,
			username: "tomjones",
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

	// seed with what we want to find
	if err := mock.InsertMultiple(db, cases[1].wantData); err != nil {
		t.Error(err)
	}

	udb := pgsql.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			user, err := udb.FindByUsername(db, tt.username)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				tt.wantData.CreatedAt = user.CreatedAt // should be set by the insert
				tt.wantData.UpdatedAt = user.UpdatedAt // should be set by the insert
				assert.Equal(t, tt.wantData, user)
			}
		})
	}
}

func TestFindByToken(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		token    string
		wantData *sandpiper.User
	}{
		{
			name:    "User does not exist",
			wantErr: true,
			token:   "notExists",
		},
		{
			name:  "Success",
			token: "loginrefresh",
			wantData: &sandpiper.User{
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
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &sandpiper.User{})

	// seed with what we want to find
	if err := mock.InsertMultiple(db, cases[1].wantData); err != nil {
		t.Error(err)
	}

	udb := pgsql.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			user, err := udb.FindByToken(db, tt.token)
			assert.Equal(t, tt.wantErr, err != nil)

			if tt.wantData != nil {
				tt.wantData.CreatedAt = user.CreatedAt // should be set by the insert
				tt.wantData.UpdatedAt = user.UpdatedAt // should be set by the insert
				assert.Equal(t, tt.wantData, user)

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
				ID:        1,
				FirstName: "Z",
				LastName:  "Freak",
				Phone:     "123456",
				Username:  "newUsername",
			},
			wantData: &sandpiper.User{
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

	// seed with original values before update
	if err := mock.InsertMultiple(db, cases[0].usr); err != nil {
		t.Error(err)
	}

	udb := pgsql.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := udb.Update(db, tt.wantData)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				user := &sandpiper.User{ID: tt.usr.ID}
				if err := db.Select(user); err != nil {
					t.Error(err)
				}
				tt.wantData.UpdatedAt = user.UpdatedAt // should be set by the update
				tt.wantData.CreatedAt = user.CreatedAt // should be set by the update
				tt.wantData.LastLogin = user.LastLogin // should be set by the update (??)
				assert.Equal(t, tt.wantData, user)
			}
		})
	}
}
