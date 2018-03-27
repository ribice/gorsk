package account_test

import (
	"context"
	"testing"

	"github.com/ribice/gorsk/internal/mock"
	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"

	"github.com/ribice/gorsk/internal/errors"

	"github.com/ribice/gorsk/internal"
	"github.com/ribice/gorsk/internal/account"
	"github.com/ribice/gorsk/internal/mock/mockdb"
)

func TestCreate(t *testing.T) {
	type args struct {
		c   *gin.Context
		req *model.User
	}
	cases := []struct {
		name     string
		args     args
		wantErr  bool
		wantData *model.User
		adb      *mockdb.Account
		udb      *mockdb.User
		rbac     *mock.RBAC
	}{{
		name: "Fail on is lower role",
		rbac: &mock.RBAC{
			AccountCreateFn: func(*gin.Context, int, int, int) bool {
				return false
			}},
		wantErr: true,
		args: args{req: &model.User{
			FirstName: "John",
			LastName:  "Doe",
			Username:  "JohnDoe",
			RoleID:    1,
			Password:  "Thranduil8822",
		}},
	},
		{
			name: "Success",
			args: args{req: &model.User{
				FirstName: "John",
				LastName:  "Doe",
				Username:  "JohnDoe",
				RoleID:    1,
				Password:  "Thranduil8822",
			}},
			adb: &mockdb.Account{
				CreateFn: func(ctx context.Context, u *model.User) error {
					u.CreatedAt = mock.TestTime(2000)
					u.UpdatedAt = mock.TestTime(2000)
					u.Base.ID = 1
					return nil
				},
			},
			rbac: &mock.RBAC{
				AccountCreateFn: func(*gin.Context, int, int, int) bool {
					return true
				}},
			wantData: &model.User{
				Base: model.Base{
					ID:        1,
					CreatedAt: mock.TestTime(2000),
					UpdatedAt: mock.TestTime(2000),
				},
				FirstName: "John",
				LastName:  "Doe",
				Username:  "JohnDoe",
				RoleID:    1,
			}}}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := account.New(tt.adb, tt.udb, tt.rbac)
			err := s.Create(tt.args.c, tt.args.req)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				tt.args.req.Password = tt.wantData.Password
				assert.Equal(t, tt.wantData, tt.args.req)
			}
		})
	}
}

func TestChangePassword(t *testing.T) {
	type args struct {
		c       *gin.Context
		oldpass string
		newpass string
		id      int
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
		udb     *mockdb.User
		adb     *mockdb.Account
		rbac    *mock.RBAC
	}{
		{
			name: "Fail on EnforceUser",
			args: args{id: 1},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c *gin.Context, id int) bool {
					return id != 1
				}},
			wantErr: true,
		},
		{
			name:    "Fail on ViewUser",
			args:    args{id: 1},
			wantErr: true,
			rbac: &mock.RBAC{
				EnforceUserFn: func(c *gin.Context, id int) bool {
					return true
				}},
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
			name: "Fail on PasswordMatch",
			args: args{id: 1, oldpass: "hunter123"},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c *gin.Context, id int) bool {
					return true
				}},
			wantErr: true,
			udb: &mockdb.User{
				ViewFn: func(c context.Context, id int) (*model.User, error) {
					return &model.User{
						Password: "IncorrectHashedPassword",
					}, nil
				},
			},
		},
		{
			name: "Success",
			args: args{id: 1, oldpass: "hunter123", newpass: "password"},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c *gin.Context, id int) bool {
					return true
				}},
			udb: &mockdb.User{
				ViewFn: func(c context.Context, id int) (*model.User, error) {
					return &model.User{
						Password: "$2a$10$udRBroNGBeOYwSWCVzf6Lulg98uAoRCIi4t75VZg84xgw6EJbFNsG",
					}, nil
				},
			},
			adb: &mockdb.Account{
				// Check whether password was hashed correctly
				ChangePasswordFn: func(c context.Context, usr *model.User) error {
					return nil
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := account.New(tt.adb, tt.udb, tt.rbac)
			err := s.ChangePassword(tt.args.c, tt.args.oldpass, tt.args.newpass, tt.args.id)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
