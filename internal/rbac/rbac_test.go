package rbac_test

import (
	"testing"

	"github.com/labstack/echo"

	"github.com/stretchr/testify/assert"

	"github.com/ribice/gorsk/internal"
	"github.com/ribice/gorsk/internal/mock"
	"github.com/ribice/gorsk/internal/rbac"
)

func TestNew(t *testing.T) {
	rbacService := rbac.New(nil)
	if rbacService == nil {
		t.Error("RBAC Service not initialized")
	}
}

func TestEnforceRole(t *testing.T) {
	type args struct {
		ctx  echo.Context
		role model.AccessRole
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Not authorized",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"role"}, int8(3)), role: model.SuperAdminRole},
			wantErr: true,
		},
		{
			name:    "Authorized",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"role"}, int8(0)), role: model.CompanyAdminRole},
			wantErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.New(nil)
			res := rbacSvc.EnforceRole(tt.args.ctx, tt.args.role)
			assert.Equal(t, tt.wantErr, res == echo.ErrForbidden)
		})
	}
}

func TestEnforceUser(t *testing.T) {
	type args struct {
		ctx echo.Context
		id  int
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Not same user, not an admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"id", "role"}, 15, int8(3)), id: 122},
			wantErr: true,
		},
		{
			name:    "Not same user, but admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"id", "role"}, 22, int8(0)), id: 44},
			wantErr: false,
		},
		{
			name:    "Same user",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"id", "role"}, 8, int8(3)), id: 8},
			wantErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.New(nil)
			res := rbacSvc.EnforceUser(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.wantErr, res == echo.ErrForbidden)
		})
	}
}

func TestEnforceCompany(t *testing.T) {
	type args struct {
		ctx echo.Context
		id  int
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Not same company, not an admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "role"}, 7, int8(5)), id: 9},
			wantErr: true,
		},
		{
			name:    "Same company, not company admin or admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "role"}, 22, int8(5)), id: 22},
			wantErr: true,
		},
		{
			name:    "Same company, company admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "role"}, 5, int8(3)), id: 5},
			wantErr: false,
		},
		{
			name:    "Not same company but admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "role"}, 8, int8(2)), id: 9},
			wantErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.New(nil)
			res := rbacSvc.EnforceCompany(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.wantErr, res == echo.ErrForbidden)
		})
	}
}

func TestEnforceLocation(t *testing.T) {
	type args struct {
		ctx echo.Context
		id  int
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Not same location, not an admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"location_id", "role"}, 7, int8(5)), id: 9},
			wantErr: true,
		},
		{
			name:    "Same location, not company admin or admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"location_id", "role"}, 22, int8(5)), id: 22},
			wantErr: true,
		},
		{
			name:    "Same location, company admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"location_id", "role"}, 5, int8(3)), id: 5},
			wantErr: false,
		},
		{
			name:    "Location admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"location_id", "role"}, 5, int8(4)), id: 5},
			wantErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.New(nil)
			res := rbacSvc.EnforceLocation(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.wantErr, res == echo.ErrForbidden)
		})
	}
}

func TestAccountCreate(t *testing.T) {
	type args struct {
		ctx         echo.Context
		roleID      int
		company_id  int
		location_id int
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Different location, company, creating user role, not an admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "location_id", "role"}, 2, 3, int8(5)), roleID: 5, company_id: 7, location_id: 8},
			wantErr: true,
		},
		{
			name:    "Same location, not company, creating user role, not an admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "location_id", "role"}, 2, 3, int8(5)), roleID: 5, company_id: 2, location_id: 8},
			wantErr: true,
		},
		{
			name:    "Different location, company, creating user role, not an admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "location_id", "role"}, 2, 3, int8(3)), roleID: 4, company_id: 2, location_id: 4},
			wantErr: false,
		},
		{
			name:    "Same location, company, creating user role, not an admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "location_id", "role"}, 2, 3, int8(3)), roleID: 5, company_id: 2, location_id: 3},
			wantErr: false,
		},
		{
			name:    "Same location, company, creating user role, admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "location_id", "role"}, 2, 3, int8(3)), roleID: 5, company_id: 2, location_id: 3},
			wantErr: false,
		},
		{
			name:    "Different everything, admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "location_id", "role"}, 2, 3, int8(1)), roleID: 2, company_id: 7, location_id: 4},
			wantErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.New(nil)
			res := rbacSvc.AccountCreate(tt.args.ctx, tt.args.roleID, tt.args.company_id, tt.args.location_id)
			assert.Equal(t, tt.wantErr, res == echo.ErrForbidden)
		})
	}
}

func TestIsLowerRole(t *testing.T) {
	ctx := mock.EchoCtxWithKeys([]string{"role"}, int8(3))
	rbacSvc := rbac.New(nil)
	if rbacSvc.IsLowerRole(ctx, model.AccessRole(4)) != nil {
		t.Error("The requested user is higher role than the user requesting it")
	}
	if rbacSvc.IsLowerRole(ctx, model.AccessRole(2)) == nil {
		t.Error("The requested user is lower role than the user requesting it")
	}
}
