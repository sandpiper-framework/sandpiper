// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package password_test

import (
	"errors"
	"testing"

	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"autocare.org/sandpiper/pkg/api/password"
	"autocare.org/sandpiper/pkg/shared/mock"
	"autocare.org/sandpiper/pkg/shared/mock/mockdb"
	"autocare.org/sandpiper/pkg/shared/model"
)

func TestChange(t *testing.T) {
	type args struct {
		oldpass string
		newpass string
		id      int
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
		udb     *mockdb.User
		rbac    *mock.RBAC
		sec     *mock.Secure
	}{
		{
			name: "Fail on EnforceUser",
			args: args{id: 1},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return errors.New("generic error")
				}},
			wantErr: true,
		},
		{
			name:    "Fail on ViewUser",
			args:    args{id: 1},
			wantErr: true,
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return nil
				}},
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, id int) (*sandpiper.User, error) {
					if id != 1 {
						return nil, nil
					}
					return nil, errors.New("generic error")
				},
			},
		},
		{
			name: "Fail on PasswordMatch",
			args: args{id: 1, oldpass: "hunter123"},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return nil
				}},
			wantErr: true,
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, id int) (*sandpiper.User, error) {
					return &sandpiper.User{
						Password: "HashedPassword",
					}, nil
				},
			},
			sec: &mock.Secure{
				HashMatchesPasswordFn: func(string, string) bool {
					return false
				},
			},
		},
		{
			name: "Fail on InsecurePassword",
			args: args{id: 1, oldpass: "hunter123"},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return nil
				}},
			wantErr: true,
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, id int) (*sandpiper.User, error) {
					return &sandpiper.User{
						Password: "HashedPassword",
					}, nil
				},
			},
			sec: &mock.Secure{
				HashMatchesPasswordFn: func(string, string) bool {
					return true
				},
				PasswordFn: func(string, ...string) bool {
					return false
				},
			},
		},
		{
			name: "Success",
			args: args{id: 1, oldpass: "hunter123", newpass: "password"},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return nil
				}},
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, id int) (*sandpiper.User, error) {
					return &sandpiper.User{
						Password: "$2a$10$udRBroNGBeOYwSWCVzf6Lulg98uAoRCIi4t75VZg84xgw6EJbFNsG",
					}, nil
				},
				UpdateFn: func(orm.DB, *sandpiper.User) error {
					return nil
				},
			},
			sec: &mock.Secure{
				HashMatchesPasswordFn: func(string, string) bool {
					return true
				},
				PasswordFn: func(string, ...string) bool {
					return true
				},
				HashFn: func(string) string {
					return "hash3d"
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := password.New(nil, tt.udb, tt.rbac, tt.sec)
			err := s.Change(nil, tt.args.id, tt.args.oldpass, tt.args.newpass)
			assert.Equal(t, tt.wantErr, err != nil)
			// Check whether password was changed
		})
	}
}

func TestInitialize(t *testing.T) {
	p := password.Initialize(nil, nil, nil)
	if p == nil {
		t.Error("password service not initialized")
	}
}
