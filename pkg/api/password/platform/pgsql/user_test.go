// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package pgsql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"sandpiper/pkg/api/password/platform/pgsql"
	"sandpiper/pkg/shared/mock"
	"sandpiper/pkg/shared/model"
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

	if err := mock.InsertMultiple(db, cases[0].usr); err != nil {
		t.Error(err)
	}

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
