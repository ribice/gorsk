package query_test

import (
	"testing"

	"github.com/labstack/echo"

	"github.com/ribice/gorsk/internal"

	"github.com/ribice/gorsk/internal/platform/query"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	type args struct {
		user *model.AuthUser
	}
	cases := []struct {
		name     string
		args     args
		wantData *model.ListQuery
		wantErr  error
	}{
		{
			name: "Super admin user",
			args: args{user: &model.AuthUser{
				Role: model.SuperAdminRole,
			}},
		},
		{
			name: "Company admin user",
			args: args{user: &model.AuthUser{
				Role:      model.CompanyAdminRole,
				CompanyID: 1,
			}},
			wantData: &model.ListQuery{
				Query: "company_id = ?",
				ID:    1},
		},
		{
			name: "Location admin user",
			args: args{user: &model.AuthUser{
				Role:       model.LocationAdminRole,
				CompanyID:  1,
				LocationID: 2,
			}},
			wantData: &model.ListQuery{
				Query: "location_id = ?",
				ID:    2},
		},
		{
			name: "Normal user",
			args: args{user: &model.AuthUser{
				Role: model.UserRole,
			}},
			wantErr: echo.ErrForbidden,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			q, err := query.List(tt.args.user)
			assert.Equal(t, tt.wantData, q)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
