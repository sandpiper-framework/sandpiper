package jwt_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"autocare.org/sandpiper/pkg/shared/middleware/jwt"
	"autocare.org/sandpiper/pkg/shared/mock"
	"autocare.org/sandpiper/pkg/shared/model"
)

func echoHandler(mw ...echo.MiddlewareFunc) *echo.Echo {
	e := echo.New()
	for _, v := range mw {
		e.Use(v)
	}
	e.GET("/hello", hwHandler)
	return e
}

func hwHandler(c echo.Context) error {
	return c.String(200, "Hello World")
}

func TestMWFunc(t *testing.T) {
	cases := []struct {
		name       string
		wantStatus int
		header     string
		signMethod string
	}{
		{
			name:       "Empty header",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Header not containing Bearer",
			header:     "notBearer",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Invalid header",
			header:     mock.HeaderInvalid(),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Header Success",
			header:     mock.HeaderValid(),
			wantStatus: http.StatusOK,
		},
	}
	jwtMW, _ := jwt.New("jwtsecret", "HS256", 60, 5)
	ts := httptest.NewServer(echoHandler(jwtMW.MWFunc()))
	defer ts.Close()
	path := ts.URL + "/hello"
	client := &http.Client{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", path, nil)
			req.Header.Set("Authorization", tt.header)
			res, err := client.Do(req)
			if err != nil {
				t.Fatal("Cannot create http request")
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestGenerateToken(t *testing.T) {
	cases := map[string]struct {
		algo         string
		secret       string
		minSecretLen int
		req          sandpiper.User
		wantErr      bool
		want         string
	}{
		"invalid algo": {
			algo:    "invalid",
			wantErr: true,
		},
		"secret not set": {
			algo:    "HS256",
			wantErr: true,
		},
		"invalid secret length": {
			algo:    "HS256",
			secret:  "123",
			wantErr: true,
		},
		"invalid secret length with min defined": {
			algo:         "HS256",
			minSecretLen: 4,
			secret:       "123",
			wantErr:      true,
		},
		"success": {
			algo:         "HS256",
			secret:       "g0r$kt3$t1ng",
			minSecretLen: 1,
			req: sandpiper.User{
				ID:        1,
				Username:  "johndoe",
				Email:     "johndoe@mail.com",
				Role:      sandpiper.SuperAdminRole,
				CompanyID: mock.TestUUID(1),
			},
			want: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			jwtSvc, err := jwt.New(tt.algo, tt.secret, 60, tt.minSecretLen)
			assert.Equal(t, tt.wantErr, err != nil)
			if err == nil && !tt.wantErr {
				token, _, _ := jwtSvc.GenerateToken(&tt.req)
				assert.Equal(t, tt.want, strings.Split(token, ".")[0])
			}
		})
	}
}
