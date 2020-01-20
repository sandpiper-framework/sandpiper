// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package transport_test

// Company test

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

	"autocare.org/sandpiper/pkg/api/company"
	"autocare.org/sandpiper/pkg/api/company/transport"
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
		wantResp   *sandpiper.Company
		sdb        *mockdb.Company
		rbac       *mock.RBAC
	}{ // todo: fix all of these req's to company ones
		{
			name:       "CREATE Fail on validation",
			req:        `{"first_name":"John","last_name":"Doe","username":"ju","password":"hunter123","password_confirm":"hunter123","email":"johndoe@gmail.com","company_id":1,"role_id":300}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "CREATE Fail on invalid role",
			req:  `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","password_confirm":"hunter123","email":"johndoe@gmail.com","company_id":1,"role_id":50}`,
			rbac: &mock.RBAC{
				AccountCreateFn: func(c echo.Context, roleID sandpiper.AccessLevel, companyID uuid.UUID) error {
					return echo.ErrForbidden
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "CREATE Fail on RBAC",
			req:  `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","password_confirm":"hunter123","email":"johndoe@gmail.com","company_id":1,"role_id":200}`,
			rbac: &mock.RBAC{
				AccountCreateFn: func(c echo.Context, roleID sandpiper.AccessLevel, companyID uuid.UUID) error {
					return echo.ErrForbidden
				},
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "CREATE Success",
			req:  `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","password_confirm":"hunter123","email":"johndoe@gmail.com","company_id":1,"role_id":200}`,
			rbac: &mock.RBAC{
				AccountCreateFn: func(c echo.Context, roleID sandpiper.AccessLevel, companyID uuid.UUID) error {
					return nil
				},
			},
			sdb: &mockdb.Company{
				CreateFn: func(db orm.DB, comp sandpiper.Company) (*sandpiper.Company, error) {
					comp.ID = mock.TestUUID(1)
					comp.CreatedAt = mock.TestTime(2018)
					comp.UpdatedAt = mock.TestTime(2018)
					return &comp, nil
				},
			},
			wantResp: &sandpiper.Company{
				ID:        mock.TestUUID(1),
				CreatedAt: mock.TestTime(2018),
				UpdatedAt: mock.TestTime(2018),
				Name:      "Acme Brakes",
				Active:    true,
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			transport.NewHTTP(company.New(nil, tt.sdb, tt.rbac, nil), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/companies"
			res, err := http.Post(path, "application/json", bytes.NewBufferString(tt.req))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(sandpiper.Company)
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
		Users []sandpiper.Company `json:"companies"`
		Page  int                 `json:"page"`
	}
	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *listResponse
		sdb        *mockdb.Company
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
			sdb: &mockdb.Company{
				ListFn: func(db orm.DB, q *sandpiper.Scope, p *sandpiper.Pagination) ([]sandpiper.Company, error) {
					if p.Limit == 100 && p.Offset == 100 {
						return []sandpiper.Company{
							{
								ID:        mock.TestUUID(1),
								CreatedAt: mock.TestTime(2001),
								UpdatedAt: mock.TestTime(2002),
								Name:      "Acme Brakes",
								Active:    true,
							},
							{
								ID:        mock.TestUUID(2),
								CreatedAt: mock.TestTime(2004),
								UpdatedAt: mock.TestTime(2005),
								Name:      "Acme Brakes",
								Active:    true,
							},
						}, nil
					}
					return nil, errors.New("generic error")
				},
			},
			wantStatus: http.StatusOK,
			wantResp: &listResponse{
				Users: []sandpiper.Company{
					{
						ID:        mock.TestUUID(1),
						CreatedAt: mock.TestTime(2001),
						UpdatedAt: mock.TestTime(2002),
						Name:      "Acme Brakes",
						Active:    true,
					},
					{
						ID:        mock.TestUUID(2),
						CreatedAt: mock.TestTime(2004),
						UpdatedAt: mock.TestTime(2005),
						Name:      "Acme Brakes",
						Active:    true,
					},
				}, Page: 1},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			transport.NewHTTP(company.New(nil, tt.sdb, tt.rbac, nil), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/companies" + tt.req
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
		wantResp   *sandpiper.Company
		sdb        *mockdb.Company
		rbac       *mock.RBAC
		sec        *mock.Secure
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
			sdb: &mockdb.Company{
				ViewFn: func(db orm.DB, id uuid.UUID) (*sandpiper.Company, error) {
					return &sandpiper.Company{
						ID:        mock.TestUUID(1),
						CreatedAt: mock.TestTime(2001),
						UpdatedAt: mock.TestTime(2002),
						Name:      "Acme Brakes",
						Active:    true,
					}, nil
				},
			},
			wantStatus: http.StatusOK,
			wantResp: &sandpiper.Company{
				ID:        mock.TestUUID(1),
				CreatedAt: mock.TestTime(2001),
				UpdatedAt: mock.TestTime(2002),
				Name:      "Acme Brakes",
				Active:    true,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			transport.NewHTTP(company.New(nil, tt.sdb, tt.rbac, tt.sec), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/companies/" + tt.req
			res, err := http.Get(path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(sandpiper.Company)
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
		wantResp   *sandpiper.Company
		sdb        *mockdb.Company
		rbac       *mock.RBAC
		sec        *mock.Secure
	}{
		{
			name:       "Invalid request",
			id:         `a`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Fail on validation",
			id:         `1`,
			req:        `{"first_name":"john","last_name":"doe","phone":"321321"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on RBAC",
			id:   `1`,
			req:  `{"first_name":"john","last_name":"doe","phone":"321321"}`,
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
			req:  `{"first_name":"john","last_name":"doe","phone":"321321"}`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(echo.Context, int) error {
					return nil
				},
			},
			sdb: &mockdb.Company{
				ViewFn: func(db orm.DB, id uuid.UUID) (*sandpiper.Company, error) {
					return &sandpiper.Company{
						ID:        mock.TestUUID(1),
						CreatedAt: mock.TestTime(2001),
						UpdatedAt: mock.TestTime(2002),
						Name:      "Acme Brakes",
						Active:    true,
					}, nil
				},
				UpdateFn: func(db orm.DB, comp *sandpiper.Company) error {
					comp.UpdatedAt = mock.TestTime(2010)
					comp.Active = false
					return nil
				},
			},
			wantStatus: http.StatusOK,
			wantResp: &sandpiper.Company{
				ID:        mock.TestUUID(1),
				CreatedAt: mock.TestTime(2001),
				UpdatedAt: mock.TestTime(2002),
				Name:      "Acme Brakes",
				Active:    false,
			},
		},
	}

	client := http.Client{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			transport.NewHTTP(company.New(nil, tt.sdb, tt.rbac, tt.sec), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/companies/" + tt.id
			req, _ := http.NewRequest("PATCH", path, bytes.NewBufferString(tt.req))
			req.Header.Set("Content-Type", "application/json")
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(sandpiper.Company)
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
		id         string
		wantStatus int
		mdb        *mockdb.Company
		rbac       *mock.RBAC
		sec        *mock.Secure
	}{
		{
			name:       "Invalid request",
			id:         `a`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on RBAC",
			id:   `1`,
			mdb: &mockdb.Company{
				ViewFn: func(db orm.DB, id uuid.UUID) (*sandpiper.Company, error) {
					return &sandpiper.Company{}, nil // todo: what should we do here? (was role:)
				},
			},
			rbac: &mock.RBAC{
				IsLowerRoleFn: func(echo.Context, sandpiper.AccessLevel) error {
					return echo.ErrForbidden
				},
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			id:   `1`,
			mdb: &mockdb.Company{
				ViewFn: func(db orm.DB, id uuid.UUID) (*sandpiper.Company, error) {
					return &sandpiper.Company{}, nil // todo: what should we do here? (was role:)
				},
				DeleteFn: func(orm.DB, *sandpiper.Company) error {
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
			transport.NewHTTP(company.New(nil, tt.mdb, tt.rbac, tt.sec), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/companies/" + tt.id
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
