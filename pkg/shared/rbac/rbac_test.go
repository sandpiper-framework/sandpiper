// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package rbac_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"sandpiper/pkg/shared/mock"
	"sandpiper/pkg/shared/model"
	"sandpiper/pkg/shared/rbac"
)

func TestUser(t *testing.T) {
	ctx := mock.EchoCtxWithKeys([]string{
		"id", "company_id", "username", "email", "role"},
		9, mock.TestUUID(1), "sandy", "sandy@gmail.com", sandpiper.SuperAdminRole)
	wantUser := &sandpiper.AuthUser{
		ID:        9,
		Username:  "sandy",
		CompanyID: mock.TestUUID(1),
		Email:     "sandy@gmail.com",
		Role:      sandpiper.SuperAdminRole,
	}
	rbacSvc := rbac.New(sandpiper.PrimaryServer)
	assert.Equal(t, wantUser, rbacSvc.CurrentUser(ctx))
}

func TestEnforceRole(t *testing.T) {
	type args struct {
		ctx  echo.Context
		role sandpiper.AccessLevel
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Not authorized",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"role"}, sandpiper.CompanyAdminRole), role: sandpiper.SuperAdminRole},
			wantErr: true,
		},
		{
			name:    "Authorized",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"role"}, sandpiper.SuperAdminRole), role: sandpiper.CompanyAdminRole},
			wantErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.New(sandpiper.PrimaryServer)
			res := rbacSvc.EnforceRole(tt.args.ctx, tt.args.role)
			assert.Equal(t, tt.wantErr, res == echo.ErrForbidden)
		})
	}
}

func TestEnforceUser(t *testing.T) {
	type args struct {
		ctx echo.Context
		id  int
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Not same user, not an admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"id", "role"}, 15, sandpiper.SyncRole), id: 122},
			wantErr: true,
		},
		{
			name:    "Not same user, but admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"id", "role"}, 22, sandpiper.SuperAdminRole), id: 44},
			wantErr: false,
		},
		{
			name:    "Same user",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"id", "role"}, 8, sandpiper.AdminRole), id: 8},
			wantErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.New(sandpiper.PrimaryServer)
			res := rbacSvc.EnforceUser(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.wantErr, res == echo.ErrForbidden)
		})
	}
}

func TestEnforceCompany(t *testing.T) {
	type args struct {
		ctx echo.Context
		id  uuid.UUID
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Not same company, not an admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "role"}, mock.TestUUID(1), sandpiper.SyncRole), id: mock.TestUUID(2)},
			wantErr: true,
		},
		{
			name:    "Same company, not company admin or admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "role"}, mock.TestUUID(1), sandpiper.SyncRole), id: mock.TestUUID(1)},
			wantErr: true,
		},
		{
			name:    "Same company, company admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "role"}, mock.TestUUID(1), sandpiper.CompanyAdminRole), id: mock.TestUUID(1)},
			wantErr: false,
		},
		{
			name:    "Not same company but admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "role"}, mock.TestUUID(1), sandpiper.AdminRole), id: mock.TestUUID(2)},
			wantErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.New(sandpiper.PrimaryServer)
			res := rbacSvc.EnforceCompany(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.wantErr, res == echo.ErrForbidden)
		})
	}
}

func TestAccountCreate(t *testing.T) {
	type args struct {
		ctx        echo.Context
		role       sandpiper.AccessLevel
		company_id uuid.UUID
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Different company, creating user role, not an admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "role"}, mock.TestUUID(1), sandpiper.SyncRole), role: 200, company_id: mock.TestUUID(2)},
			wantErr: true,
		},
		{
			name:    "Different company, creating user role, not an admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "role"}, mock.TestUUID(1), sandpiper.CompanyAdminRole), role: 120, company_id: mock.TestUUID(2)},
			wantErr: false,
		},
		{
			name:    "Same company, creating user role, not an admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "role"}, mock.TestUUID(1), sandpiper.CompanyAdminRole), role: 120, company_id: mock.TestUUID(1)},
			wantErr: false,
		},
		{
			name:    "Same company, creating user role, admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "role"}, mock.TestUUID(1), sandpiper.CompanyAdminRole), role: 120, company_id: mock.TestUUID(1)},
			wantErr: false,
		},
		{
			name:    "Different everything, admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "role"}, mock.TestUUID(1), sandpiper.AdminRole), role: 110, company_id: mock.TestUUID(2)},
			wantErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.New(sandpiper.PrimaryServer)
			res := rbacSvc.AccountCreate(tt.args.ctx, tt.args.role, tt.args.company_id)
			assert.Equal(t, tt.wantErr, res == echo.ErrForbidden)
		})
	}
}

func TestIsLowerRole(t *testing.T) {
	ctx := mock.EchoCtxWithKeys([]string{"role"}, sandpiper.CompanyAdminRole)
	rbacSvc := rbac.New(sandpiper.PrimaryServer)
	if rbacSvc.IsLowerRole(ctx, sandpiper.AdminRole) == nil {
		t.Error("The requested user is lower role than the user requesting it")
	}
}
