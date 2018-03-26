package request_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ribice/gorsk/cmd/api/request"
)

func TestAccountCreate(t *testing.T) {
	type errResp struct {
		wantStatus int
		wantResp   string
	}
	cases := []struct {
		name     string
		e        *errResp
		req      string
		wantErr  bool
		wantData *request.Register
	}{
		{
			name:    "Fail on binding JSON",
			wantErr: true,
			req:     `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","password_confirm":"hunter1234","email":"johndoe@gmail.com","company_id":1,"location_id":2}`,
			e: &errResp{
				wantStatus: http.StatusBadRequest,
				wantResp:   `{"message":["RoleID is required, but was not received"]}`,
			},
		},
		{
			name:    "Fail on password match",
			wantErr: true,
			req:     `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","password_confirm":"hunter1234","email":"johndoe@gmail.com","company_id":1,"location_id":2,"role_id":3}`,
			e: &errResp{
				wantStatus: http.StatusBadRequest,
				wantResp:   `{"message":"passwords do not match"}`,
			},
		},
		{
			name:    "Fail on non-existent role_id",
			wantErr: true,
			req:     `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","password_confirm":"hunter123","email":"johndoe@gmail.com","company_id":1,"location_id":2,"role_id":9}`,
			e: &errResp{
				wantStatus: http.StatusBadRequest,
			},
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
			gin.SetMode(gin.TestMode)
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("POST", "", bytes.NewBufferString(tt.req))
			reg, err := request.AccountCreate(c)
			if tt.e != nil {
				if tt.e.wantStatus != w.Code {
					t.Errorf("Expected status %v, received %v", tt.e.wantStatus, w.Code)
				}
				if tt.e.wantResp != "" && tt.e.wantResp != w.Body.String() {
					t.Errorf("Expected response %v, received %v", tt.e.wantResp, w.Body.String())
				}
			}
			if !reflect.DeepEqual(tt.wantData, reg) {
				t.Errorf("Expected %v, received %v", tt.wantData, reg)
			}
			if tt.wantErr != (err != nil) {
				t.Errorf("Expected err = %v, but was %v", tt.wantErr, err != nil)
			}
			if tt.wantErr != c.IsAborted() {
				t.Error("Expected context to be aborted but was not")
			}
		})
	}
}

func TestPasswordChange(t *testing.T) {
	type errResp struct {
		wantStatus int
		wantResp   string
	}
	cases := []struct {
		name     string
		e        *errResp
		id       string
		req      string
		wantErr  bool
		wantData *request.Password
	}{
		{
			name:    "Fail on ID param",
			wantErr: true,
			id:      "NaN",
			e: &errResp{
				wantStatus: http.StatusBadRequest,
			},
		},
		{
			name:    "Fail on binding JSON",
			wantErr: true,
			id:      "1",
			e: &errResp{
				wantStatus: http.StatusBadRequest,
				wantResp:   `{"message":["NewPasswordConfirm is required, but was not received"]}`,
			},
			req: `{"new_password":"new_password","old_password":"my_old_password"}`,
		},
		{
			name:    "Not matching passwords",
			wantErr: true,
			id:      "1",
			e: &errResp{
				wantStatus: http.StatusBadRequest,
				wantResp:   `{"message":"passwords do not match"}`,
			},
			req: `{"new_password":"new_password","old_password":"my_old_password", "new_password_confirm":"new_password_cf"}`,
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
	gin.SetMode(gin.TestMode)
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{gin.Param{Key: "id", Value: tt.id}}
			if tt.req != "" {
				c.Request, _ = http.NewRequest("POST", "", bytes.NewBufferString(tt.req))
			}
			pw, err := request.PasswordChange(c)
			if tt.e != nil {
				if tt.e.wantStatus != w.Code {
					t.Errorf("Expected status %v, received %v", tt.e.wantStatus, w.Code)
				}
				if tt.e.wantResp != "" && tt.e.wantResp != w.Body.String() {
					t.Errorf("Expected response %v, received %v", tt.e.wantResp, w.Body.String())
				}
			}
			if !reflect.DeepEqual(tt.wantData, pw) {
				t.Errorf("Expected %v, received %v", tt.wantData, pw)
			}
			if tt.wantErr != (err != nil) {
				t.Errorf("Expected err = %v, but was %v", tt.wantErr, err != nil)
			}
			if tt.wantErr != c.IsAborted() {
				t.Error("Expected context to be aborted but was not")
			}
		})
	}
}
