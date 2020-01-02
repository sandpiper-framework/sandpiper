package transport_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"autocare.org/sandpiper/pkg/internal/middleware/jwt"
	"autocare.org/sandpiper/pkg/internal/model"
	"autocare.org/sandpiper/pkg/internal/server"
	"autocare.org/sandpiper/pkg/api/auth"
	"autocare.org/sandpiper/pkg/api/auth/transport"
	"autocare.org/sandpiper/test/mock"
	"autocare.org/sandpiper/test/mock/mockdb"
)

func TestLogin(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *sandpiper.AuthToken
		udb        *mockdb.User
		jwt        *mock.JWT
		sec        *mock.Secure
	}{
		{
			name:       "Invalid request",
			req:        `{"username":"juzernejm"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Fail on FindByUsername",
			req:        `{"username":"juzernejm","password":"hunter123"}`,
			wantStatus: http.StatusInternalServerError,
			udb: &mockdb.User{
				FindByUsernameFn: func(orm.DB, string) (*sandpiper.User, error) {
					return nil, errors.New("generic error")
				},
			},
		},
		{
			name:       "Success",
			req:        `{"username":"juzernejm","password":"hunter123"}`,
			wantStatus: http.StatusOK,
			udb: &mockdb.User{
				FindByUsernameFn: func(orm.DB, string) (*sandpiper.User, error) {
					return &sandpiper.User{
						Password: "hunter123",
						Active:   true,
					}, nil
				},
				UpdateFn: func(db orm.DB, u *sandpiper.User) error {
					return nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(*sandpiper.User) (string, string, error) {
					return "jwttokenstring", mock.TestTime(2018).Format(time.RFC3339), nil
				},
			},
			sec: &mock.Secure{
				HashMatchesPasswordFn: func(string, string) bool {
					return true
				},
				TokenFn: func(string) string {
					return "refreshtoken"
				},
			},
			wantResp: &sandpiper.AuthToken{Token: "jwttokenstring", Expires: mock.TestTime(2018).Format(time.RFC3339), RefreshToken: "refreshtoken"},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			transport.NewHTTP(auth.New(nil, tt.udb, tt.jwt, tt.sec, nil), r, nil)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/login"
			res, err := http.Post(path, "application/json", bytes.NewBufferString(tt.req))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(sandpiper.AuthToken)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				tt.wantResp.RefreshToken = response.RefreshToken
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestRefresh(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *sandpiper.RefreshToken
		udb        *mockdb.User
		jwt        *mock.JWT
	}{
		{
			name:       "Fail on FindByToken",
			req:        "refreshtoken",
			wantStatus: http.StatusInternalServerError,
			udb: &mockdb.User{
				FindByTokenFn: func(orm.DB, string) (*sandpiper.User, error) {
					return nil, errors.New("generic error")
				},
			},
		},
		{
			name:       "Success",
			req:        "refreshtoken",
			wantStatus: http.StatusOK,
			udb: &mockdb.User{
				FindByTokenFn: func(orm.DB, string) (*sandpiper.User, error) {
					return &sandpiper.User{
						Username: "johndoe",
						Active:   true,
					}, nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(*sandpiper.User) (string, string, error) {
					return "jwttokenstring", mock.TestTime(2018).Format(time.RFC3339), nil
				},
			},
			wantResp: &sandpiper.RefreshToken{Token: "jwttokenstring", Expires: mock.TestTime(2018).Format(time.RFC3339)},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			transport.NewHTTP(auth.New(nil, tt.udb, tt.jwt, nil, nil), r, nil)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/refresh/" + tt.req
			res, err := http.Get(path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(sandpiper.RefreshToken)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestMe(t *testing.T) {
	cases := []struct {
		name       string
		wantStatus int
		wantResp   *sandpiper.User
		header     string
		udb        *mockdb.User
		rbac       *mock.RBAC
	}{
		{
			name:       "Fail on user view",
			wantStatus: http.StatusInternalServerError,
			udb: &mockdb.User{
				ViewFn: func(orm.DB, int) (*sandpiper.User, error) {
					return nil, errors.New("generic error")
				},
			},
			rbac: &mock.RBAC{
				CurrentUserFn: func(echo.Context) *sandpiper.AuthUser {
					return &sandpiper.AuthUser{ID: 1}
				},
			},
			header: mock.HeaderValid(),
		},
		{
			name:       "Success",
			wantStatus: http.StatusOK,
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, i int) (*sandpiper.User, error) {
					return &sandpiper.User{
						ID:        i,
						CompanyID: mock.TestUUID(1),
						Email:     "john@mail.com",
						FirstName: "John",
						LastName:  "Doe",
					}, nil
				},
			},
			rbac: &mock.RBAC{
				CurrentUserFn: func(echo.Context) *sandpiper.AuthUser {
					return &sandpiper.AuthUser{ID: 1}
				},
			},
			header: mock.HeaderValid(),
			wantResp: &sandpiper.User{
				ID:        1,
				CompanyID: mock.TestUUID(1),
				Email:     "john@mail.com",
				FirstName: "John",
				LastName:  "Doe",
			},
		},
	}

	client := &http.Client{}
	jwtMW := jwt.New("jwtsecret", "HS256", 60)

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			transport.NewHTTP(auth.New(nil, tt.udb, nil, nil, tt.rbac), r, jwtMW.MWFunc())
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/me"
			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Authorization", tt.header)
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(sandpiper.User)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}
