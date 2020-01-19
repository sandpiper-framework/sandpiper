// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package transport_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"autocare.org/sandpiper/pkg/api/slice"
	"autocare.org/sandpiper/pkg/api/slice/transport"
	"autocare.org/sandpiper/pkg/shared/mock"
	"autocare.org/sandpiper/pkg/shared/mock/mockdb"
	"autocare.org/sandpiper/pkg/shared/model"
	"autocare.org/sandpiper/pkg/shared/server"
)

func TestCreate(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *sandpiper.Slice
		udb        *mockdb.Slice
		rbac       *mock.RBAC
	}{
		{
			name:       "Fail on validation",
			req:        `{"first_name":"John","last_name":"Doe","username":"ju","password":"hunter123","password_confirm":"hunter123","email":"johndoe@gmail.com","company_id": "10000000-0000-0000-0000-000000000000","role":300}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Fail on non-matching passwords",
			req:        `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","password_confirm":"hunter1234","email":"johndoe@gmail.com","company_id": "10000000-0000-0000-0000-000000000000","role":300}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on invalid role",
			req:  `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","password_confirm":"hunter123","email":"johndoe@gmail.com","company_id": "10000000-0000-0000-0000-000000000000","role":50}`,
			rbac: &mock.RBAC{
				AccountCreateFn: func(c echo.Context, roleID sandpiper.AccessLevel, companyID uuid.UUID) error {
					return echo.ErrForbidden
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on RBAC",
			req:  `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","password_confirm":"hunter123","email":"johndoe@gmail.com","company_id":"10000000-0000-0000-0000-000000000000","role":200}`,
			rbac: &mock.RBAC{
				AccountCreateFn: func(c echo.Context, roleID sandpiper.AccessLevel, companyID uuid.UUID) error {
					return echo.ErrForbidden
				},
			},
			wantStatus: http.StatusForbidden,
		},

		{
			name: "Success",
			req:  `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","password_confirm":"hunter123","email":"johndoe@gmail.com","company_id":"10000000-0000-0000-0000-000000000000","role":200}`,
			rbac: &mock.RBAC{
				AccountCreateFn: func(c echo.Context, roleID sandpiper.AccessLevel, companyID uuid.UUID) error {
					return nil
				},
			},
			udb: &mockdb.Slice{
				CreateFn: func(db orm.DB, slice sandpiper.Slice) (*sandpiper.Slice, error) {
					slice.ID = mock.TestUUID(1)
					slice.CreatedAt = mock.TestTime(2018)
					slice.UpdatedAt = mock.TestTime(2018)
					return &slice, nil
				},
			},
			wantResp: &sandpiper.Slice{
				ID:        mock.TestUUID(1),
				Name:      "AAP Brakes",
				CreatedAt: mock.TestTime(2018),
				UpdatedAt: mock.TestTime(2018),
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			transport.NewHTTP(slice.New(nil, tt.udb, tt.rbac, nil), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/slices"
			res, err := http.Post(path, "application/json", bytes.NewBufferString(tt.req))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(sandpiper.Slice)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestList(t *testing.T) {
	type listResponse struct {
		Slices []sandpiper.Slice `json:"slices"`
		Page   int               `json:"page"`
	}
	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *listResponse
		udb        *mockdb.Slice
		rbac       *mock.RBAC
	}{
		{
			name:       "Invalid request",
			req:        `?limit=2222&page=-1`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on query list",
			req:  `?limit=100&page=1`,
			rbac: &mock.RBAC{
				CurrentUserFn: func(c echo.Context) *sandpiper.AuthUser {
					return &sandpiper.AuthUser{
						ID:        1,
						CompanyID: mock.TestUUID(1),
						Role:      sandpiper.CompanyAdminRole,
						Email:     "john@mail.com",
					}
				}},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			req:  `?limit=100&page=1`,
			rbac: &mock.RBAC{
				CurrentUserFn: func(c echo.Context) *sandpiper.AuthUser {
					return &sandpiper.AuthUser{
						ID:        1,
						CompanyID: mock.TestUUID(1),
						Role:      sandpiper.SuperAdminRole,
						Email:     "john@mail.com",
					}
				}},
			udb: &mockdb.Slice{
				ListFn: func(db orm.DB, q *sandpiper.Clause, p *sandpiper.Pagination) ([]sandpiper.Slice, error) {
					if p.Limit == 100 && p.Offset == 100 {
						return []sandpiper.Slice{
							{
								ID:        mock.TestUUID(1),
								Name:      "AAP Brakes",
								CreatedAt: mock.TestTime(2001),
								UpdatedAt: mock.TestTime(2002),
							},
							{
								ID:        mock.TestUUID(2),
								Name:      "AAP Wipers",
								CreatedAt: mock.TestTime(2004),
								UpdatedAt: mock.TestTime(2005),
							},
						}, nil
					}
					return nil, errors.New("generic error")
				},
			},
			wantStatus: http.StatusOK,
			wantResp: &listResponse{
				Slices: []sandpiper.Slice{
					{
						ID:        mock.TestUUID(1),
						Name:      "AAP Brakes",
						CreatedAt: mock.TestTime(2001),
						UpdatedAt: mock.TestTime(2002),
					},
					{
						ID:        mock.TestUUID(2),
						Name:      "AAP Wipers",
						CreatedAt: mock.TestTime(2004),
						UpdatedAt: mock.TestTime(2005),
					},
				}, Page: 1},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			transport.NewHTTP(slice.New(nil, tt.udb, tt.rbac, nil), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/slices" + tt.req
			res, err := http.Get(path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(listResponse)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestView(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *sandpiper.Slice
		udb        *mockdb.Slice
		rbac       *mock.RBAC
	}{
		{
			name:       "Invalid request",
			req:        `a`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on RBAC",
			req:  `1`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(echo.Context, int) error {
					return echo.ErrForbidden
				},
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			req:  `1`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(echo.Context, int) error {
					return nil
				},
			},
			udb: &mockdb.Slice{
				ViewFn: func(db orm.DB, id uuid.UUID) (*sandpiper.Slice, error) {
					return &sandpiper.Slice{
						ID:        mock.TestUUID(1),
						Name:      "AAP Brakes",
						CreatedAt: mock.TestTime(2000),
						UpdatedAt: mock.TestTime(2000),
					}, nil
				},
			},
			wantStatus: http.StatusOK,
			wantResp: &sandpiper.Slice{
				ID:        mock.TestUUID(1),
				Name:      "AAP Brakes",
				CreatedAt: mock.TestTime(2000),
				UpdatedAt: mock.TestTime(2000),
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			transport.NewHTTP(slice.New(nil, tt.udb, tt.rbac, nil), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/slices/" + tt.req
			res, err := http.Get(path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(sandpiper.Slice)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestUpdate(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		id         string
		wantStatus int
		wantResp   *sandpiper.Slice
		udb        *mockdb.Slice
		rbac       *mock.RBAC
	}{
		{
			name:       "Invalid request",
			id:         `a`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Fail on validation",
			id:         `1`,
			req:        `{"first_name":"j","last_name":"okocha","mobile":"123456","phone":"321321","address":"home"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on RBAC",
			id:   `1`,
			req:  `{"first_name":"jj","last_name":"okocha","mobile":"123456","phone":"321321","address":"home"}`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(echo.Context, int) error {
					return echo.ErrForbidden
				},
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			id:   `1`,
			req:  `{"first_name":"jj","last_name":"okocha","phone":"321321","address":"home"}`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(echo.Context, int) error {
					return nil
				},
			},
			udb: &mockdb.Slice{
				ViewFn: func(db orm.DB, id uuid.UUID) (*sandpiper.Slice, error) {
					return &sandpiper.Slice{
						ID:           mock.TestUUID(1),
						Name:         "AAP Brakes",
						ContentCount: 1,
						CreatedAt:    mock.TestTime(2000),
						UpdatedAt:    mock.TestTime(2000),
					}, nil
				},
				UpdateFn: func(db orm.DB, slice *sandpiper.Slice) error {
					slice.UpdatedAt = mock.TestTime(2010)
					slice.ContentCount = 2
					return nil
				},
			},
			wantStatus: http.StatusOK,
			wantResp: &sandpiper.Slice{
				ID:           mock.TestUUID(1),
				Name:         "AAP Brakes",
				ContentCount: 1,
				CreatedAt:    mock.TestTime(2000),
				UpdatedAt:    mock.TestTime(2000),
			},
		},
	}

	client := http.Client{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			transport.NewHTTP(slice.New(nil, tt.udb, tt.rbac, nil), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/slices/" + tt.id
			req, _ := http.NewRequest("PATCH", path, bytes.NewBufferString(tt.req))
			req.Header.Set("Content-Type", "application/json")
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(sandpiper.Slice)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		name       string
		id         string // to allow testing a bad request
		cid        uuid.UUID
		wantStatus int
		mdb        *mockdb.Slice
		rbac       *mock.RBAC
	}{
		{
			name:       "Invalid request",
			id:         "123",
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on RBAC",
			id:   mock.TestUUID(1).String(),
			cid:  mock.TestUUID(1),
			mdb: &mockdb.Slice{
				ViewBySubFn: func(db orm.DB, cid uuid.UUID, id uuid.UUID) (*sandpiper.Slice, error) {
					return &sandpiper.Slice{
						Name: "AAP Brakes",
					}, nil
				},
			},
			rbac: &mock.RBAC{
				IsLowerRoleFn: func(echo.Context, sandpiper.AccessLevel) error {
					return echo.ErrForbidden
				},
			},
			wantStatus: http.StatusForbidden,
		},
		{ // todo: this test looks wrong... (should allow sysadmin or company-admin of subscribed companies)
			name: "Success",
			id:   mock.TestUUID(1).String(),
			cid:  mock.TestUUID(1),
			mdb: &mockdb.Slice{
				ViewBySubFn: func(db orm.DB, cid uuid.UUID, id uuid.UUID) (*sandpiper.Slice, error) {
					return &sandpiper.Slice{
						Name: "AAP Brakes",
					}, nil
				},
				DeleteFn: func(orm.DB, *sandpiper.Slice) error {
					return nil
				},
			},
			rbac: &mock.RBAC{
				IsLowerRoleFn: func(echo.Context, sandpiper.AccessLevel) error {
					return nil
				},
			},
			wantStatus: http.StatusOK,
		},
	}

	client := http.Client{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			transport.NewHTTP(slice.New(nil, tt.mdb, tt.rbac, nil), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/slices/" + tt.id
			req, _ := http.NewRequest("DELETE", path, nil)
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}
