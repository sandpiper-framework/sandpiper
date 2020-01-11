// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package scope_test

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"autocare.org/sandpiper/pkg/internal/mock"
	"autocare.org/sandpiper/pkg/internal/model"
	"autocare.org/sandpiper/pkg/internal/scope"
)

func TestList(t *testing.T) {
	type args struct {
		user *sandpiper.AuthUser
	}
	cases := []struct {
		name     string
		args     args
		wantData *scope.Clause
		wantErr  error
	}{
		{
			name: "Super admin user",
			args: args{user: &sandpiper.AuthUser{
				Role: sandpiper.SuperAdminRole,
			}},
		},
		{
			name: "Company admin user",
			args: args{user: &sandpiper.AuthUser{
				Role:      sandpiper.CompanyAdminRole,
				CompanyID: mock.TestUUID(1),
			}},
			wantData: &scope.Clause{
				Condition: "company_id = ?",
				ID:        mock.TestUUID(1)},
		},
		{
			name: "Normal user",
			args: args{user: &sandpiper.AuthUser{
				Role: sandpiper.UserRole,
			}},
			wantErr: echo.ErrForbidden,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			q, err := scope.Limit(tt.args.user)
			assert.Equal(t, tt.wantData, q)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
