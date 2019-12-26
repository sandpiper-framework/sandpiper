package slice_test

import (
	"errors"
	"testing"

	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"autocare.org/sandpiper/internal/model"
	"autocare.org/sandpiper/pkg/api/user"
	"autocare.org/sandpiper/test/mock"
	"autocare.org/sandpiper/test/mock/mockdb"
)

func TestCreate(t *testing.T) {
	type args struct {
		c   echo.Context
		req sandpiper.User
	}
	cases := []struct {
		name     string
		args     args
		wantErr  bool
		wantData *sandpiper.User
		udb      *mockdb.User
		rbac     *mock.RBAC
		sec      *mock.Secure
	}{{
		name: "Fail on is lower role",
		rbac: &mock.RBAC{
			AccountCreateFn: func(echo.Context, sandpiper.AccessRole, int, int) error {
				return errors.New("generic error")
			}},
		wantErr: true,
		args: args{req: sandpiper.User{
			FirstName: "John",
			LastName:  "Doe",
			Username:  "JohnDoe",
			RoleID:    1,
			Password:  "Thranduil8822",
		}},
	},
		{
			name: "Success",
			args: args{req: sandpiper.User{
				FirstName: "John",
				LastName:  "Doe",
				Username:  "JohnDoe",
				RoleID:    1,
				Password:  "Thranduil8822",
			}},
			udb: &mockdb.User{
				CreateFn: func(db orm.DB, u sandpiper.User) (*sandpiper.User, error) {
					u.CreatedAt = mock.TestTime(2000)
					u.UpdatedAt = mock.TestTime(2000)
					u.Base.ID = 1
					return &u, nil
				},
			},
			rbac: &mock.RBAC{
				AccountCreateFn: func(echo.Context, sandpiper.AccessRole, int, int) error {
					return nil
				}},
			sec: &mock.Secure{
				HashFn: func(string) string {
					return "h4$h3d"
				},
			},
			wantData: &sandpiper.User{
				Base: sandpiper.Base{
					ID:        1,
					CreatedAt: mock.TestTime(2000),
					UpdatedAt: mock.TestTime(2000),
				},
				FirstName: "John",
				LastName:  "Doe",
				Username:  "JohnDoe",
				RoleID:    1,
				Password:  "h4$h3d",
			}}}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := user.New(nil, tt.udb, tt.rbac, tt.sec)
			usr, err := s.Create(tt.args.c, tt.args.req)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantData, usr)
		})
	}
}

func TestView(t *testing.T) {
	type args struct {
		c  echo.Context
		id int
	}
	cases := []struct {
		name     string
		args     args
		wantData *sandpiper.User
		wantErr  error
		udb      *mockdb.User
		rbac     *mock.RBAC
	}{
		{
			name: "Fail on RBAC",
			args: args{id: 5},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return errors.New("generic error")
				}},
			wantErr: errors.New("generic error"),
		},
		{
			name: "Success",
			args: args{id: 1},
			wantData: &sandpiper.User{
				Base: sandpiper.Base{
					ID:        1,
					CreatedAt: mock.TestTime(2000),
					UpdatedAt: mock.TestTime(2000),
				},
				FirstName: "John",
				LastName:  "Doe",
				Username:  "JohnDoe",
			},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return nil
				}},
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, id int) (*sandpiper.User, error) {
					if id == 1 {
						return &sandpiper.User{
							Base: sandpiper.Base{
								ID:        1,
								CreatedAt: mock.TestTime(2000),
								UpdatedAt: mock.TestTime(2000),
							},
							FirstName: "John",
							LastName:  "Doe",
							Username:  "JohnDoe",
						}, nil
					}
					return nil, nil
				}},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := user.New(nil, tt.udb, tt.rbac, nil)
			usr, err := s.View(tt.args.c, tt.args.id)
			assert.Equal(t, tt.wantData, usr)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestList(t *testing.T) {
	type args struct {
		c   echo.Context
		pgn *sandpiper.Pagination
	}
	cases := []struct {
		name     string
		args     args
		wantData []sandpiper.User
		wantErr  bool
		udb      *mockdb.User
		rbac     *mock.RBAC
	}{
		{
			name: "Fail on query List",
			args: args{c: nil, pgn: &sandpiper.Pagination{
				Limit:  100,
				Offset: 200,
			}},
			wantErr: true,
			rbac: &mock.RBAC{
				UserFn: func(c echo.Context) *sandpiper.AuthUser {
					return &sandpiper.AuthUser{
						ID:         1,
						CompanyID:  2,
						LocationID: 3,
						Role:       sandpiper.UserRole,
					}
				}}},
		{
			name: "Success",
			args: args{c: nil, pgn: &sandpiper.Pagination{
				Limit:  100,
				Offset: 200,
			}},
			rbac: &mock.RBAC{
				UserFn: func(c echo.Context) *sandpiper.AuthUser {
					return &sandpiper.AuthUser{
						ID:         1,
						CompanyID:  2,
						LocationID: 3,
						Role:       sandpiper.AdminRole,
					}
				}},
			udb: &mockdb.User{
				ListFn: func(orm.DB, *sandpiper.ListQuery, *sandpiper.Pagination) ([]sandpiper.User, error) {
					return []sandpiper.User{
						{
							Base: sandpiper.Base{
								ID:        1,
								CreatedAt: mock.TestTime(1999),
								UpdatedAt: mock.TestTime(2000),
							},
							FirstName: "John",
							LastName:  "Doe",
							Email:     "johndoe@gmail.com",
							Username:  "johndoe",
						},
						{
							Base: sandpiper.Base{
								ID:        2,
								CreatedAt: mock.TestTime(2001),
								UpdatedAt: mock.TestTime(2002),
							},
							FirstName: "Hunter",
							LastName:  "Logan",
							Email:     "logan@aol.com",
							Username:  "hunterlogan",
						},
					}, nil
				}},
			wantData: []sandpiper.User{
				{
					Base: sandpiper.Base{
						ID:        1,
						CreatedAt: mock.TestTime(1999),
						UpdatedAt: mock.TestTime(2000),
					},
					FirstName: "John",
					LastName:  "Doe",
					Email:     "johndoe@gmail.com",
					Username:  "johndoe",
				},
				{
					Base: sandpiper.Base{
						ID:        2,
						CreatedAt: mock.TestTime(2001),
						UpdatedAt: mock.TestTime(2002),
					},
					FirstName: "Hunter",
					LastName:  "Logan",
					Email:     "logan@aol.com",
					Username:  "hunterlogan",
				}},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := user.New(nil, tt.udb, tt.rbac, nil)
			usrs, err := s.List(tt.args.c, tt.args.pgn)
			assert.Equal(t, tt.wantData, usrs)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}

}

func TestDelete(t *testing.T) {
	type args struct {
		c  echo.Context
		id int
	}
	cases := []struct {
		name    string
		args    args
		wantErr error
		udb     *mockdb.User
		rbac    *mock.RBAC
	}{
		{
			name:    "Fail on ViewUser",
			args:    args{id: 1},
			wantErr: errors.New("generic error"),
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
			name: "Fail on RBAC",
			args: args{id: 1},
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, id int) (*sandpiper.User, error) {
					return &sandpiper.User{
						Base: sandpiper.Base{
							ID:        id,
							CreatedAt: mock.TestTime(1999),
							UpdatedAt: mock.TestTime(2000),
						},
						FirstName: "John",
						LastName:  "Doe",
						Role: &sandpiper.Role{
							AccessLevel: sandpiper.UserRole,
						},
					}, nil
				},
			},
			rbac: &mock.RBAC{
				IsLowerRoleFn: func(echo.Context, sandpiper.AccessRole) error {
					return errors.New("generic error")
				}},
			wantErr: errors.New("generic error"),
		},
		{
			name: "Success",
			args: args{id: 1},
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, id int) (*sandpiper.User, error) {
					return &sandpiper.User{
						Base: sandpiper.Base{
							ID:        id,
							CreatedAt: mock.TestTime(1999),
							UpdatedAt: mock.TestTime(2000),
						},
						FirstName: "John",
						LastName:  "Doe",
						Role: &sandpiper.Role{
							AccessLevel: sandpiper.AdminRole,
							ID:          2,
							Name:        "Admin",
						},
					}, nil
				},
				DeleteFn: func(db orm.DB, usr *sandpiper.User) error {
					return nil
				},
			},
			rbac: &mock.RBAC{
				IsLowerRoleFn: func(echo.Context, sandpiper.AccessRole) error {
					return nil
				}},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := user.New(nil, tt.udb, tt.rbac, nil)
			err := s.Delete(tt.args.c, tt.args.id)
			if err != tt.wantErr {
				t.Errorf("Expected error %v, received %v", tt.wantErr, err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type args struct {
		c   echo.Context
		upd *user.Update
	}
	cases := []struct {
		name     string
		args     args
		wantData *sandpiper.User
		wantErr  error
		udb      *mockdb.User
		rbac     *mock.RBAC
	}{
		{
			name: "Fail on RBAC",
			args: args{upd: &user.Update{
				ID: 1,
			}},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return errors.New("generic error")
				}},
			wantErr: errors.New("generic error"),
		},
		{
			name: "Fail on Update",
			args: args{upd: &user.Update{
				ID: 1,
			}},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return nil
				}},
			wantErr: errors.New("generic error"),
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, id int) (*sandpiper.User, error) {
					return &sandpiper.User{
						Base: sandpiper.Base{
							ID:        1,
							CreatedAt: mock.TestTime(1990),
							UpdatedAt: mock.TestTime(1991),
						},
						CompanyID:  1,
						LocationID: 2,
						RoleID:     3,
						FirstName:  "John",
						LastName:   "Doe",
						Mobile:     "123456",
						Phone:      "234567",
						Address:    "Work Address",
						Email:      "golang@go.org",
					}, nil
				},
				UpdateFn: func(db orm.DB, usr *sandpiper.User) error {
					return errors.New("generic error")
				},
			},
		},
		{
			name: "Success",
			args: args{upd: &user.Update{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Mobile:    "123456",
				Phone:     "234567",
			}},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return nil
				}},
			wantData: &sandpiper.User{
				Base: sandpiper.Base{
					ID:        1,
					CreatedAt: mock.TestTime(1990),
					UpdatedAt: mock.TestTime(2000),
				},
				CompanyID:  1,
				LocationID: 2,
				RoleID:     3,
				FirstName:  "John",
				LastName:   "Doe",
				Mobile:     "123456",
				Phone:      "234567",
				Address:    "Work Address",
				Email:      "golang@go.org",
			},
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, id int) (*sandpiper.User, error) {
					return &sandpiper.User{
						Base: sandpiper.Base{
							ID:        1,
							CreatedAt: mock.TestTime(1990),
							UpdatedAt: mock.TestTime(2000),
						},
						CompanyID:  1,
						LocationID: 2,
						RoleID:     3,
						FirstName:  "John",
						LastName:   "Doe",
						Mobile:     "123456",
						Phone:      "234567",
						Address:    "Work Address",
						Email:      "golang@go.org",
					}, nil
				},
				UpdateFn: func(db orm.DB, usr *sandpiper.User) error {
					usr.UpdatedAt = mock.TestTime(2000)
					return nil
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := user.New(nil, tt.udb, tt.rbac, nil)
			usr, err := s.Update(tt.args.c, tt.args.upd)
			assert.Equal(t, tt.wantData, usr)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestInitialize(t *testing.T) {
	u := user.Initialize(nil, nil, nil)
	if u == nil {
		t.Error("User service not initialized")
	}
}
