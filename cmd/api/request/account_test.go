package request_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ribice/gorsk/cmd/api/request"
	"github.com/ribice/gorsk/internal/mock"
)

func TestAccountCreate(t *testing.T) {
	cases := []struct {
		name     string
		req      string
		wantErr  bool
		wantData *request.Register
	}{
		{
			name:    "Fail on validating JSON",
			wantErr: true,
			req:     `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","password_confirm":"hunter1234","email":"johndoe@gmail.com","company_id":1,"location_id":2}`,
		},
		{
			name:    "Fail on password match",
			wantErr: true,
			req:     `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","password_confirm":"hunter1234","email":"johndoe@gmail.com","company_id":1,"location_id":2,"role_id":3}`,
		},
		{
			name:    "Fail on non-existent role_id",
			wantErr: true,
			req:     `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","password_confirm":"hunter123","email":"johndoe@gmail.com","company_id":1,"location_id":2,"role_id":9}`,
		},
		{
			name: "Success",
			req:  `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","password_confirm":"hunter123","email":"johndoe@gmail.com","company_id":1,"location_id":2,"role_id":2}`,
			wantData: &request.Register{
				FirstName:       "John",
				LastName:        "Doe",
				Username:        "juzernejm",
				Password:        "hunter123",
				PasswordConfirm: "hunter123",
				Email:           "johndoe@gmail.com",
				CompanyID:       1,
				LocationID:      2,
				RoleID:          2,
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "", bytes.NewBufferString(tt.req))
			c := mock.EchoCtx(req, w)
			reg, err := request.AccountCreate(c)
			assert.Equal(t, tt.wantData, reg)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestPasswordChange(t *testing.T) {
	cases := []struct {
		name     string
		id       string
		req      string
		wantErr  bool
		wantData *request.Password
	}{
		{
			name:    "Fail on ID param",
			wantErr: true,
			id:      "NaN",
		},
		{
			name:    "Fail on binding JSON",
			wantErr: true,
			id:      "1",
			req:     `{"new_password":"new_password","old_password":"my_old_password"}`,
		},
		{
			name:    "Not matching passwords",
			wantErr: true,
			id:      "1",
			req:     `{"new_password":"new_password","old_password":"my_old_password", "new_password_confirm":"new_password_cf"}`,
		},
		{
			name: "Success",
			id:   "10",
			req:  `{"new_password":"newpassw","old_password":"oldpassw", "new_password_confirm":"newpassw"}`,
			wantData: &request.Password{
				ID:                 10,
				NewPassword:        "newpassw",
				NewPasswordConfirm: "newpassw",
				OldPassword:        "oldpassw",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(tt.req))
			c := mock.EchoCtx(req, w)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			pw, err := request.PasswordChange(c)
			assert.Equal(t, tt.wantData, pw)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
