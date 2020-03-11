// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package transport_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"autocare.org/sandpiper/pkg/primary/password"
	"autocare.org/sandpiper/pkg/primary/password/transport"
	"autocare.org/sandpiper/pkg/shared/mock"
	"autocare.org/sandpiper/pkg/shared/mock/mockdb"
	"autocare.org/sandpiper/pkg/shared/model"
	"autocare.org/sandpiper/pkg/shared/server"
)

func TestChangePassword(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		wantStatus int
		id         string
		udb        *mockdb.User
		rbac       *mock.RBAC
		sec        *mock.Secure
	}{
		{
			name:       "NaN",
			wantStatus: http.StatusBadRequest,
			id:         "abc",
		},
		{
			name:       "Fail on Bind",
			req:        `{"new_password":"new","old_password":"my_old_password", "new_password_confirm":"new"}`,
			wantStatus: http.StatusBadRequest,
			id:         "1",
		},
		{
			name:       "Different passwords",
			req:        `{"new_password":"new_password","old_password":"my_old_password", "new_password_confirm":"new_password_cf"}`,
			wantStatus: http.StatusBadRequest,
			id:         "1",
		},
		{
			name: "Fail on RBAC",
			req:  `{"new_password":"newpassw","old_password":"oldpassw", "new_password_confirm":"newpassw"}`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return echo.ErrForbidden
				},
			},
			id:         "1",
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			req:  `{"new_password":"newpassw","old_password":"oldpassw", "new_password_confirm":"newpassw"}`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return nil
				},
			},
			id: "1",
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, id int) (*sandpiper.User, error) {
					return &sandpiper.User{
						Password: "oldPassword",
					}, nil
				},
				UpdateFn: func(db orm.DB, usr *sandpiper.User) error {
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
					return "hashedPassword"
				},
			},
			wantStatus: http.StatusOK,
		},
	}

	client := &http.Client{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			transport.NewHTTP(password.New(nil, tt.udb, tt.rbac, tt.sec), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/password/" + tt.id
			req, err := http.NewRequest("PATCH", path, bytes.NewBufferString(tt.req))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}