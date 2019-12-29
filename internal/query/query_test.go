package query_test

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"autocare.org/sandpiper/internal/model"
	"autocare.org/sandpiper/internal/query"
	"autocare.org/sandpiper/test/mock"
)

func TestList(t *testing.T) {
	type args struct {
		user *sandpiper.AuthUser
	}
	cases := []struct {
		name     string
		args     args
		wantData *sandpiper.ListQuery
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
			wantData: &sandpiper.ListQuery{
				Query: "company_id = ?",
				ID:    mock.TestUUID(1)},
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
			q, err := query.List(tt.args.user)
			assert.Equal(t, tt.wantData, q)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

