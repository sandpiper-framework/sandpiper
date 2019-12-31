// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package slice_test

import (
	"errors"
	"testing"
	"time"

	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"autocare.org/sandpiper/internal/model"
	"autocare.org/sandpiper/pkg/api/slice"
	"autocare.org/sandpiper/test/mock"
	"autocare.org/sandpiper/test/mock/mockdb"
)

func TestCreate(t *testing.T) {
	type args struct {
		ctx echo.Context
		req sandpiper.Slice
	}
	cases := []struct {
		name     string
		args     args
		wantErr  bool
		wantData *sandpiper.Slice
		mdb      *mockdb.Slice
		rbac     *mock.RBAC
	}{{
		name: "CREATE Fails as Standard User",
		rbac: &mock.RBAC{
			EnforceRoleFn: func(echo.Context, sandpiper.AccessLevel) error {
				return errors.New("forbidden error")
			}},
		wantErr: true,
		args: args{
			ctx: mock.EchoCtxWithKeys([]string{"role"}, sandpiper.UserRole),
			req: sandpiper.Slice{
				Name:   "AAP Brake Friction",
				ContentHash: "4468e5deabf5e6d0740cd1a77df56f67093ec943",
				ContentCount: 1,
				LastUpdate: time.Now(),
			}},
	},
		{
			name: "CREATE Succeeds as Company Admin",
			args: args{req: sandpiper.Slice{
				Name:   "AAP Brake Friction",
				ContentHash: "4468e5deabf5e6d0740cd1a77df56f67093ec943",
				ContentCount: 1,
				LastUpdate: time.Now(),
			}},
			mdb: &mockdb.Slice{
				CreateFn: func(db orm.DB, u sandpiper.Slice) (*sandpiper.Slice, error) {
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
			wantData: &sandpiper.Slice{
				ID:        mock.TestUUID(1),
				Name:   "AAP Brake Friction",
				ContentHash: "4468e5deabf5e6d0740cd1a77df56f67093ec943",
				ContentCount: 1,
				LastUpdate: time.Now(),
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := slice.New(nil, tt.mdb, tt.rbac, nil)
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
		wantData *sandpiper.Slice
		wantErr  error
		mdb      *mockdb.Slice
		rbac     *mock.RBAC
	}{
		{
			name: "VIEW Fails with User Permissions",
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
			wantData: &sandpiper.Slice{
				ID:        mock.TestUUID(1),
				CreatedAt: mock.TestTime(2000),
				UpdatedAt: mock.TestTime(2000),
				Name:   "AAP Brake Friction",
				ContentHash: "4468e5deabf5e6d0740cd1a77df56f67093ec943",
				ContentCount: 1,
				LastUpdate: time.Now(),
			},
			rbac: &mock.RBAC{
				EnforceCompanyFn: func(c echo.Context, id uuid.UUID) error {
					return nil
				}},
			mdb: &mockdb.Slice{
				ViewFn: func(db orm.DB, id uuid.UUID) (*sandpiper.Slice, error) {
					if id == mock.TestUUID(1) {
						return &sandpiper.Slice{
							ID:        mock.TestUUID(1),
							CreatedAt: mock.TestTime(2000),
							UpdatedAt: mock.TestTime(2000),
							Name:   "AAP Brake Friction",
							ContentHash: "4468e5deabf5e6d0740cd1a77df56f67093ec943",
							ContentCount: 1,
							LastUpdate: time.Now(),
						}, nil
					}
					return nil, nil
				}},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := slice.New(nil, tt.mdb, tt.rbac, nil)
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
		wantData []sandpiper.Slice
		wantErr  bool
		mdb      *mockdb.Slice
		rbac     *mock.RBAC
	}{
		{
			name: "LIST Failed on query List",
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
			name: "LIST Succeeded",
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
			mdb: &mockdb.Slice{
				ListFn: func(orm.DB, *scope.Clause, *sandpiper.Pagination) ([]sandpiper.Slice, error) {
					return []sandpiper.Slice{
						{
							ID:        mock.TestUUID(1),
							CreatedAt: mock.TestTime(1999),
							UpdatedAt: mock.TestTime(2000),
							Name:   "AAP Brake Friction",
							ContentHash: "4468e5deabf5e6d0740cd1a77df56f67093ec943",
							ContentCount: 1,
							LastUpdate: time.Now(),
						},
						{
							ID:        mock.TestUUID(2),
							CreatedAt: mock.TestTime(2001),
							UpdatedAt: mock.TestTime(2002),
							Name:   "AAP Premium Wipers",
							ContentHash: "da39a3ee5e6b4b0d3255bfef95601890afd80709",
							ContentCount: 1,
							LastUpdate: time.Now(),
						},
					}, nil
				}},
			wantData: []sandpiper.Slice{
				{
					ID:        mock.TestUUID(1),
					CreatedAt: mock.TestTime(1999),
					UpdatedAt: mock.TestTime(2000),
					ContentHash: "4468e5deabf5e6d0740cd1a77df56f67093ec943",
					ContentCount: 1,
					LastUpdate: time.Now(),
				},
				{
					ID:        mock.TestUUID(2),
					CreatedAt: mock.TestTime(2001),
					UpdatedAt: mock.TestTime(2002),
					Name:   "AAP Premium Wipers",
					ContentHash: "da39a3ee5e6b4b0d3255bfef95601890afd80709",
					ContentCount: 1,
					LastUpdate: time.Now(),
				}},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := slice.New(nil, tt.mdb, tt.rbac, nil)
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
		mdb     *mockdb.Slice
		rbac    *mock.RBAC
	}{
		{
			name:    "DELETE Fail on ViewUser",
			args:    args{id: mock.TestUUID(1)},
			wantErr: errors.New("generic error"),
			mdb: &mockdb.Slice{
				ViewFn: func(db orm.DB, id uuid.UUID) (*sandpiper.Slice, error) {
					if id != mock.TestUUID(1) {
						return nil, nil
					}
					return nil, errors.New("generic error")
				},
			},
		},
		{
			name: "DELETE Fail on RBAC",
			args: args{id: mock.TestUUID(1)},
			mdb: &mockdb.Slice{
				ViewFn: func(db orm.DB, id uuid.UUID) (*sandpiper.Slice, error) {
					return &sandpiper.Slice{
						ID:        mock.TestUUID(1),
						CreatedAt: mock.TestTime(1999),
						UpdatedAt: mock.TestTime(2000),
						Name:   "AAP Brake Friction",
						ContentHash: "4468e5deabf5e6d0740cd1a77df56f67093ec943",
						ContentCount: 1,
						LastUpdate: time.Now(),
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
			name: "DELETE Successful",
			args: args{id: mock.TestUUID(1)},
			mdb: &mockdb.Slice{
				ViewFn: func(db orm.DB, id uuid.UUID) (*sandpiper.Slice, error) {
					return &sandpiper.Slice{
						ID:        mock.TestUUID(1),
						CreatedAt: mock.TestTime(1999),
						UpdatedAt: mock.TestTime(2000),
						Name:   "AAP Brake Friction",
						ContentHash: "4468e5deabf5e6d0740cd1a77df56f67093ec943",
						ContentCount: 1,
						LastUpdate: time.Now(),
					}, nil
				},
				DeleteFn: func(db orm.DB, usr *sandpiper.Slice) error {
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
			s := slice.New(nil, tt.mdb, tt.rbac, nil)
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
		upd *slice.Update
	}
	cases := []struct {
		name     string
		args     args
		wantData *sandpiper.Slice
		wantErr  error
		mdb      *mockdb.Slice
		rbac     *mock.RBAC
	}{
		{
			name: "UPDATE Fail on RBAC",
			args: args{upd: &slice.Update{
				ID: mock.TestUUID(1),
			}},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return errors.New("generic error")
				}},
			wantErr: errors.New("generic error"),
		},
		{
			name: "UPDATE Fail on Update",
			args: args{upd: &slice.Update{
				ID: mock.TestUUID(1),
			}},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return nil
				}},
			wantErr: errors.New("generic error"),
			mdb: &mockdb.Slice{
				ViewFn: func(db orm.DB, id uuid.UUID) (*sandpiper.Slice, error) {
					return &sandpiper.Slice{
						ID:        mock.TestUUID(1),
						CreatedAt: mock.TestTime(1990),
						UpdatedAt: mock.TestTime(1991),
						Name:   "AAP Brake Friction",
						ContentHash: "4468e5deabf5e6d0740cd1a77df56f67093ec943",
						ContentCount: 1,
						LastUpdate: time.Now(),
					}, nil
				},
				UpdateFn: func(db orm.DB, usr *sandpiper.Slice) error {
					return errors.New("generic error")
				},
			},
		},
		{
			name: "UPDATE Success",
			args: args{upd: &slice.Update{
				ID:     mock.TestUUID(1),
				Name:   "AAP Brake Friction",
				ContentHash: "4468e5deabf5e6d0740cd1a77df56f67093ec943",
				ContentCount: 1,
				LastUpdate: time.Now(),
			}},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return nil
				}},
			wantData: &sandpiper.Slice{
				ID:        mock.TestUUID(1),
				CreatedAt: mock.TestTime(1990),
				UpdatedAt: mock.TestTime(2000),
				Name:   "AAP Brake Friction",
				ContentHash: "4468e5deabf5e6d0740cd1a77df56f67093ec943",
				ContentCount: 1,
				LastUpdate: time.Now(),
			},
			mdb: &mockdb.Slice{
				ViewFn: func(db orm.DB, id uuid.UUID) (*sandpiper.Slice, error) {
					return &sandpiper.Slice{
						ID:        mock.TestUUID(1),
						CreatedAt: mock.TestTime(1990),
						UpdatedAt: mock.TestTime(2000),
						Name:   "AAP Brake Friction",
						ContentHash: "4468e5deabf5e6d0740cd1a77df56f67093ec943",
						ContentCount: 1,
						LastUpdate: time.Now(),
					}, nil
				},
				UpdateFn: func(db orm.DB, slice *sandpiper.Slice) error {
					slice.UpdatedAt = mock.TestTime(2000)
					return nil
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := slice.New(nil, tt.mdb, tt.rbac, nil)
			usr, err := s.Update(tt.args.c, tt.args.upd)
			assert.Equal(t, tt.wantData, usr)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestInitialize(t *testing.T) {
	s := slice.Initialize(nil, nil, nil)
	if s == nil {
		t.Error("Slice service not initialized")
	}
}
