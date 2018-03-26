package request_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ribice/gorsk/internal/mock"

	"github.com/gin-gonic/gin"
	"github.com/ribice/gorsk/cmd/api/request"
)

func TestUserUpdate(t *testing.T) {
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
		wantData *request.UpdateUser
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
				wantResp:   `{"message":["FirstName's value or length is less than allowed"]}`,
			},
			req: `{"first_name":"j","last_name":"okocha"}`,
		},
		{
			name: "Success",
			id:   "1",
			req:  `{"first_name":"jj","last_name":"okocha","mobile":"123456","phone":"321321","address":"home"}`,
			wantData: &request.UpdateUser{
				ID:        1,
				FirstName: mock.Str2Ptr("jj"),
				LastName:  mock.Str2Ptr("okocha"),
				Mobile:    mock.Str2Ptr("123456"),
				Phone:     mock.Str2Ptr("321321"),
				Address:   mock.Str2Ptr("home"),
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
			resp, err := request.UserUpdate(c)
			if tt.e != nil {
				if tt.e.wantStatus != w.Code {
					t.Errorf("Expected status %v, received %v", tt.e.wantStatus, w.Code)
				}
				if tt.e.wantResp != "" && tt.e.wantResp != w.Body.String() {
					t.Errorf("Expected response %v, received %v", tt.e.wantResp, w.Body.String())
				}
			}
			if !reflect.DeepEqual(tt.wantData, resp) {
				t.Errorf("Expected %v, received %v", tt.wantData, resp)
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
