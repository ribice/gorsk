package query_test

import (
	"reflect"
	"testing"

	"github.com/ribice/gorsk/internal"
	"github.com/ribice/gorsk/internal/errors"
	"github.com/ribice/gorsk/internal/platform/query"
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
			wantErr: apperr.Forbidden,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			q, err := query.List(tt.args.user)
			if !reflect.DeepEqual(tt.wantData, q) {
				t.Errorf("Expected and returned data does not match")
			}
			if err != tt.wantErr {
				t.Errorf("Expected error %v, received %v", tt.wantErr, err)
			}
		})
	}
}
