package user_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ribice/gorsk/internal/errors"

	"github.com/ribice/gorsk/internal"
	"github.com/ribice/gorsk/internal/mock"
	"github.com/ribice/gorsk/internal/mock/mockdb"
	"github.com/ribice/gorsk/internal/user"
)

func TestView(t *testing.T) {
	type args struct {
		c  *gin.Context
		id int
	}
	cases := []struct {
		name     string
		args     args
		wantData *model.User
		wantErr  error
		udb      *mockdb.User
		rbac     *mock.RBAC
	}{
		{
			name: "Fail on RBAC",
			args: args{id: 5},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c *gin.Context, id int) bool {
					return id == 1
				}},
			wantErr: apperr.Forbidden,
		},
		{
			name: "Success",
			args: args{id: 1},
			wantData: &model.User{
				Base: model.Base{
					ID:        1,
					CreatedAt: mock.TestTime(2000),
					UpdatedAt: mock.TestTime(2000),
				},
				FirstName: "John",
				LastName:  "Doe",
				Username:  "JohnDoe",
			},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c *gin.Context, id int) bool {
					return true
				}},
			udb: &mockdb.User{
				ViewFn: func(ctx context.Context, id int) (*model.User, error) {
					if id == 1 {
						return &model.User{
							Base: model.Base{
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
			s := user.New(tt.udb, tt.rbac, nil)
			usr, err := s.View(tt.args.c, tt.args.id)
			if !reflect.DeepEqual(tt.wantData, usr) {
				t.Errorf("Expected and returned data does not match")
			}
			if err != tt.wantErr {
				t.Errorf("Expected error %v, received %v", tt.wantErr, err)
			}
		})
	}
}

func TestList(t *testing.T) {
	type args struct {
		c   *gin.Context
		pgn *model.Pagination
	}
	cases := []struct {
		name     string
		args     args
		wantData []model.User
		wantErr  bool
		udb      *mockdb.User
		auth     *mock.Auth
	}{
		{
			name: "Fail on query List",
			args: args{c: &gin.Context{}, pgn: &model.Pagination{
				Limit:  100,
				Offset: 200,
			}},
			wantErr: true,
			auth: &mock.Auth{
				UserFn: func(c *gin.Context) *model.AuthUser {
					return &model.AuthUser{
						ID:         1,
						CompanyID:  2,
						LocationID: 3,
						Role:       model.UserRole,
					}
				}}},
		{
			name: "Success",
			args: args{c: &gin.Context{}, pgn: &model.Pagination{
				Limit:  100,
				Offset: 200,
			}},
			auth: &mock.Auth{
				UserFn: func(c *gin.Context) *model.AuthUser {
					return &model.AuthUser{
						ID:         1,
						CompanyID:  2,
						LocationID: 3,
						Role:       model.AdminRole,
					}
				}},
			udb: &mockdb.User{
				ListFn: func(context.Context, *model.ListQuery, *model.Pagination) ([]model.User, error) {
					return []model.User{
						{
							Base: model.Base{
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
							Base: model.Base{
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
			wantData: []model.User{
				{
					Base: model.Base{
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
					Base: model.Base{
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
			s := user.New(tt.udb, nil, tt.auth)
			usrs, err := s.List(tt.args.c, tt.args.pgn)
			if !reflect.DeepEqual(tt.wantData, usrs) {
				t.Errorf("Expected and returned data does not match")
			}
			if tt.wantErr != (err != nil) {
				t.Errorf("Expected error %v, received %v", tt.wantErr, err)
			}
		})
	}

}

func TestDelete(t *testing.T) {
	type args struct {
		c  *gin.Context
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
			wantErr: apperr.DB,
			udb: &mockdb.User{
				ViewFn: func(c context.Context, id int) (*model.User, error) {
					if id != 1 {
						return nil, nil
					}
					return nil, apperr.DB
				},
			},
		},
		{
			name: "Fail on RBAC",
			args: args{id: 1},
			udb: &mockdb.User{
				ViewFn: func(c context.Context, id int) (*model.User, error) {
					return &model.User{
						Base: model.Base{
							ID:        id,
							CreatedAt: mock.TestTime(1999),
							UpdatedAt: mock.TestTime(2000),
						},
						FirstName: "John",
						LastName:  "Doe",
						Role: &model.Role{
							AccessLevel: model.UserRole,
						},
					}, nil
				},
			},
			rbac: &mock.RBAC{
				IsLowerRoleFn: func(*gin.Context, model.AccessRole) bool {
					return false
				}},
			wantErr: apperr.Forbidden,
		},
		{
			name: "Success",
			args: args{id: 1},
			udb: &mockdb.User{
				ViewFn: func(c context.Context, id int) (*model.User, error) {
					return &model.User{
						Base: model.Base{
							ID:        id,
							CreatedAt: mock.TestTime(1999),
							UpdatedAt: mock.TestTime(2000),
						},
						FirstName: "John",
						LastName:  "Doe",
						Role: &model.Role{
							AccessLevel: model.AdminRole,
							ID:          2,
							Name:        "Admin",
						},
					}, nil
				},
				DeleteFn: func(c context.Context, usr *model.User) error {
					return nil
				},
			},
			rbac: &mock.RBAC{
				IsLowerRoleFn: func(*gin.Context, model.AccessRole) bool {
					return true
				}},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := user.New(tt.udb, tt.rbac, nil)
			err := s.Delete(tt.args.c, tt.args.id)
			if err != tt.wantErr {
				t.Errorf("Expected error %v, received %v", tt.wantErr, err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type args struct {
		c   *gin.Context
		upd *user.Update
	}
	cases := []struct {
		name     string
		args     args
		wantData *model.User
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
				EnforceUserFn: func(c *gin.Context, id int) bool {
					return false
				}},
			wantErr: apperr.Forbidden,
		},
		{
			name: "Fail on ViewUser",
			args: args{upd: &user.Update{
				ID: 1,
			}},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c *gin.Context, id int) bool {
					return true
				}},
			wantErr: apperr.DB,
			udb: &mockdb.User{
				ViewFn: func(c context.Context, id int) (*model.User, error) {
					if id != 1 {
						return nil, nil
					}
					return nil, apperr.DB
				},
			},
		},
		{
			name: "Success",
			args: args{upd: &user.Update{
				ID:        1,
				FirstName: mock.Str2Ptr("John"),
				LastName:  mock.Str2Ptr("Doe"),
				Mobile:    mock.Str2Ptr("123456"),
				Phone:     mock.Str2Ptr("234567"),
			}},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c *gin.Context, id int) bool {
					return true
				}},
			wantData: &model.User{
				Base: model.Base{
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
				ViewFn: func(c context.Context, id int) (*model.User, error) {
					if id == 1 {
						return &model.User{
							Base: model.Base{
								ID:        1,
								CreatedAt: mock.TestTime(1990),
								UpdatedAt: mock.TestTime(1991),
							},
							CompanyID:  1,
							LocationID: 2,
							RoleID:     3,
							FirstName:  "Joanna",
							LastName:   "Doep",
							Mobile:     "334455",
							Phone:      "444555",
							Address:    "Work Address",
							Email:      "golang@go.org",
						}, nil
					}
					return nil, apperr.DB
				},
				UpdateFn: func(c context.Context, usr *model.User) (*model.User, error) {
					usr.UpdatedAt = mock.TestTime(2000)
					return usr, nil
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := user.New(tt.udb, tt.rbac, nil)
			usr, err := s.Update(tt.args.c, tt.args.upd)
			if !reflect.DeepEqual(tt.wantData, usr) {
				t.Errorf("Expected and returned data does not match")
			}
			if err != tt.wantErr {
				t.Errorf("Expected error %v, received %v", tt.wantErr, err)
			}
		})
	}
}
