// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package subscription_test

import (
	"errors"
	"testing"

	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"autocare.org/sandpiper/pkg/primary/company"
	"autocare.org/sandpiper/pkg/shared/mock"
	"autocare.org/sandpiper/pkg/shared/mock/mockdb"
	"autocare.org/sandpiper/pkg/shared/model"
)

func TestCreate(t *testing.T) {
	type args struct {
		ctx echo.Context
		req sandpiper.Company
	}
	cases := []struct {
		name     string
		args     args
		wantErr  bool
		wantData *sandpiper.Company
		mdb      *mockdb.Company
		rbac     *mock.RBAC
	}{{
		name: "Fails as Standard User",
		rbac: &mock.RBAC{
			EnforceRoleFn: func(echo.Context, sandpiper.AccessLevel) error {
				return errors.New("forbidden error")
			}},
		wantErr: true,
		args: args{
			ctx: mock.EchoCtxWithKeys([]string{"role"}, sandpiper.UserRole),
			req: sandpiper.Company{
				Name:   "Acme Brakes",
				Active: true,
			}},
	},
		{
			name:    "Succeeds as Company Admin",
			wantErr: false,
			args: args{
				ctx: mock.EchoCtxWithKeys([]string{"role"}, sandpiper.CompanyAdminRole),
				req: sandpiper.Company{
					Name:   "Acme Brakes",
					Active: true,
				}},
			mdb: &mockdb.Company{
				CreateFn: func(db orm.DB, u sandpiper.Company) (*sandpiper.Company, error) {
					u.CreatedAt = mock.TestTime(2000)
					u.UpdatedAt = mock.TestTime(2000)
					u.ID = mock.TestUUID(1)
					return &u, nil
				},
			},
			rbac: &mock.RBAC{
				EnforceRoleFn: func(echo.Context, sandpiper.AccessLevel) error {
					return errors.New("forbidden error")
				}},
			wantData: &sandpiper.Company{
				ID:        mock.TestUUID(1),
				CreatedAt: mock.TestTime(2000),
				UpdatedAt: mock.TestTime(2000),
				Name:      "Acme Brakes",
				Active:    true,
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := company.New(nil, tt.mdb, tt.rbac, nil)
			res, err := s.Create(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantData, res)
		})
	}
}

func TestView(t *testing.T) {
	type args struct {
		c  echo.Context
		id uuid.UUID
	}
	cases := []struct {
		name     string
		args     args
		wantData *sandpiper.Company
		wantErr  error
		mdb      *mockdb.Company
		rbac     *mock.RBAC
	}{
		{
			name: "Fails with User Permissions",
			args: args{
				mock.EchoCtxWithKeys([]string{"role"}, sandpiper.UserRole),
				mock.TestUUID(1),
			},
			rbac: &mock.RBAC{
				EnforceCompanyFn: func(c echo.Context, id uuid.UUID) error {
					return errors.New("forbidden error")
				}},
			wantErr: errors.New("forbidden error"),
		},
		{
			name: "VIEW Success",
			args: args{
				mock.EchoCtxWithKeys([]string{"role"}, sandpiper.CompanyAdminRole),
				mock.TestUUID(1),
			},
			wantData: &sandpiper.Company{
				ID:        mock.TestUUID(1),
				CreatedAt: mock.TestTime(2000),
				UpdatedAt: mock.TestTime(2000),
				Name:      "Acme Brakes",
				Active:    true,
			},
			rbac: &mock.RBAC{
				EnforceCompanyFn: func(c echo.Context, id uuid.UUID) error {
					return nil
				}},
			mdb: &mockdb.Company{
				ViewFn: func(db orm.DB, id uuid.UUID) (*sandpiper.Company, error) {
					if id == mock.TestUUID(1) {
						return &sandpiper.Company{
							ID:        mock.TestUUID(1),
							CreatedAt: mock.TestTime(2000),
							UpdatedAt: mock.TestTime(2000),
							Name:      "Acme Brakes",
							Active:    true,
						}, nil
					}
					return nil, nil
				}},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := company.New(nil, tt.mdb, tt.rbac, nil)
			res, err := s.View(tt.args.c, tt.args.id)
			assert.Equal(t, tt.wantData, res)
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
		wantData []sandpiper.Company
		wantErr  bool
		mdb      *mockdb.Company
		rbac     *mock.RBAC
	}{
		{
			name: "Failed on query List",
			args: args{c: nil, pgn: &sandpiper.Pagination{
				Limit:  100,
				Offset: 200,
			}},
			wantErr: true,
			rbac: &mock.RBAC{
				CurrentUserFn: func(c echo.Context) *sandpiper.AuthUser {
					return &sandpiper.AuthUser{
						ID:        1,
						CompanyID: mock.TestUUID(1),
						Role:      sandpiper.CompanyAdminRole,
					}
				}}},
		{
			name: "Succeeded",
			args: args{c: nil, pgn: &sandpiper.Pagination{
				Limit:  100,
				Offset: 200,
			}},
			rbac: &mock.RBAC{
				CurrentUserFn: func(c echo.Context) *sandpiper.AuthUser {
					return &sandpiper.AuthUser{
						ID:        1,
						CompanyID: mock.TestUUID(1),
						Role:      sandpiper.AdminRole,
					}
				}},
			mdb: &mockdb.Company{
				ListFn: func(orm.DB, *sandpiper.Scope, *sandpiper.Pagination) ([]sandpiper.Company, error) {
					return []sandpiper.Company{
						{
							ID:        mock.TestUUID(1),
							CreatedAt: mock.TestTime(1999),
							UpdatedAt: mock.TestTime(2000),
							Name:      "Acme Brakes",
							Active:    true,
						},
						{
							ID:        mock.TestUUID(2),
							CreatedAt: mock.TestTime(2001),
							UpdatedAt: mock.TestTime(2002),
							Name:      "Acme Wipers",
							Active:    true,
						},
					}, nil
				}},
			wantData: []sandpiper.Company{
				{
					ID:        mock.TestUUID(1),
					CreatedAt: mock.TestTime(1999),
					UpdatedAt: mock.TestTime(2000),
					Name:      "Acme Brakes",
					Active:    true,
				},
				{
					ID:        mock.TestUUID(2),
					CreatedAt: mock.TestTime(2001),
					UpdatedAt: mock.TestTime(2002),
					Name:      "Acme Wipers",
					Active:    true,
				}},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := company.New(nil, tt.mdb, tt.rbac, nil)
			res, err := s.List(tt.args.c, tt.args.pgn)
			assert.Equal(t, tt.wantData, res)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}

}

func TestDelete(t *testing.T) {
	type args struct {
		c  echo.Context
		id uuid.UUID
	}
	cases := []struct {
		name    string
		args    args
		wantErr error
		mdb     *mockdb.Company
		rbac    *mock.RBAC
	}{
		{
			name:    "Fail on ViewUser",
			args:    args{id: mock.TestUUID(1)},
			wantErr: errors.New("generic error"),
			mdb: &mockdb.Company{
				ViewFn: func(db orm.DB, id uuid.UUID) (*sandpiper.Company, error) {
					if id != mock.TestUUID(1) {
						return nil, nil
					}
					return nil, errors.New("generic error")
				},
			},
		},
		{
			name: "Fail on RBAC",
			args: args{id: mock.TestUUID(1)},
			mdb: &mockdb.Company{
				ViewFn: func(db orm.DB, id uuid.UUID) (*sandpiper.Company, error) {
					return &sandpiper.Company{
						ID:        mock.TestUUID(1),
						CreatedAt: mock.TestTime(1999),
						UpdatedAt: mock.TestTime(2000),
						Name:      "Acme Brakes",
						Active:    true,
					}, nil
				},
			},
			rbac: &mock.RBAC{
				IsLowerRoleFn: func(echo.Context, sandpiper.AccessLevel) error {
					return errors.New("generic error")
				}},
			wantErr: errors.New("generic error"),
		},
		{
			name: "Successful",
			args: args{id: mock.TestUUID(1)},
			mdb: &mockdb.Company{
				ViewFn: func(db orm.DB, id uuid.UUID) (*sandpiper.Company, error) {
					return &sandpiper.Company{
						ID:        mock.TestUUID(1),
						CreatedAt: mock.TestTime(1999),
						UpdatedAt: mock.TestTime(2000),
						Name:      "Acme Brakes",
						Active:    true,
					}, nil
				},
				DeleteFn: func(db orm.DB, usr *sandpiper.Company) error {
					return nil
				},
			},
			rbac: &mock.RBAC{
				IsLowerRoleFn: func(echo.Context, sandpiper.AccessLevel) error {
					return nil
				}},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := company.New(nil, tt.mdb, tt.rbac, nil)
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
		upd *company.Update
	}
	cases := []struct {
		name     string
		args     args
		wantData *sandpiper.Company
		wantErr  error
		mdb      *mockdb.Company
		rbac     *mock.RBAC
	}{
		{
			name: "Fail on RBAC",
			args: args{upd: &company.Update{
				ID: mock.TestUUID(1),
			}},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return errors.New("generic error")
				}},
			wantErr: errors.New("generic error"),
		},
		{
			name: "Fail on Update",
			args: args{upd: &company.Update{
				ID: mock.TestUUID(1),
			}},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return nil
				}},
			wantErr: errors.New("generic error"),
			mdb: &mockdb.Company{
				ViewFn: func(db orm.DB, id uuid.UUID) (*sandpiper.Company, error) {
					return &sandpiper.Company{
						ID:        mock.TestUUID(1),
						CreatedAt: mock.TestTime(1990),
						UpdatedAt: mock.TestTime(1991),
						Name:      "Acme Brakes",
						Active:    true,
					}, nil
				},
				UpdateFn: func(db orm.DB, usr *sandpiper.Company) error {
					return errors.New("generic error")
				},
			},
		},
		{
			name: "Success",
			args: args{upd: &company.Update{
				ID:     mock.TestUUID(1),
				Name:   "Acme Brakes",
				Active: true,
			}},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return nil
				}},
			wantData: &sandpiper.Company{
				ID:        mock.TestUUID(1),
				CreatedAt: mock.TestTime(1990),
				UpdatedAt: mock.TestTime(2000),
				Name:      "Acme Brakes",
				Active:    true,
			},
			mdb: &mockdb.Company{
				ViewFn: func(db orm.DB, id uuid.UUID) (*sandpiper.Company, error) {
					return &sandpiper.Company{
						ID:        mock.TestUUID(1),
						CreatedAt: mock.TestTime(1990),
						UpdatedAt: mock.TestTime(2000),
						Name:      "Acme Brakes",
						Active:    true,
					}, nil
				},
				UpdateFn: func(db orm.DB, usr *sandpiper.Company) error {
					usr.UpdatedAt = mock.TestTime(2000)
					return nil
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := company.New(nil, tt.mdb, tt.rbac, nil)
			usr, err := s.Update(tt.args.c, tt.args.upd)
			assert.Equal(t, tt.wantData, usr)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestInitialize(t *testing.T) {
	s := company.Initialize(nil, nil, nil)
	if s == nil {
		t.Error("User service not initialized")
	}
}
