package query_test

import (
	"testing"

	"github.com/labstack/echo"

	"github.com/ribice/gorsk"

	"github.com/stretchr/testify/assert"

	"github.com/ribice/gorsk/pkg/utl/query"
)

func TestList(t *testing.T) {
	type args struct {
		user gorsk.AuthUser
	}
	cases := []struct {
		name     string
		args     args
		wantData *gorsk.ListQuery
		wantErr  error
	}{
		{
			name: "Super admin user",
			args: args{user: gorsk.AuthUser{
				Role: gorsk.SuperAdminRole,
			}},
		},
		{
			name: "Company admin user",
			args: args{user: gorsk.AuthUser{
				Role:      gorsk.CompanyAdminRole,
				CompanyID: 1,
			}},
			wantData: &gorsk.ListQuery{
				Query: "company_id = ?",
				ID:    1},
		},
		{
			name: "Location admin user",
			args: args{user: gorsk.AuthUser{
				Role:       gorsk.LocationAdminRole,
				CompanyID:  1,
				LocationID: 2,
			}},
			wantData: &gorsk.ListQuery{
				Query: "location_id = ?",
				ID:    2},
		},
		{
			name: "Normal user",
			args: args{user: gorsk.AuthUser{
				Role: gorsk.UserRole,
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
