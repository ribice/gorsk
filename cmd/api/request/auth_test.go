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

func TestLogin(t *testing.T) {
	type errResp struct {
		wantStatus int
		wantResp   string
	}
	cases := []struct {
		name     string
		e        *errResp
		req      string
		wantErr  bool
		wantData *request.Credentials
	}{
		{
			name:    "Fail on binding JSON",
			wantErr: true,
			req:     `{"username":"juzernejm"}`,
			e: &errResp{
				wantStatus: http.StatusBadRequest,
				wantResp:   `{"message":["Password is required, but was not received"]}`,
			},
		},
		{
			name: "Success",
			req:  `{"username":"juzernejm","password":"hunter123"}`,
			wantData: &request.Credentials{
				Username: "juzernejm",
				Password: "hunter123",
			},
		},
	}
	gin.SetMode(gin.TestMode)
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("POST", "", bytes.NewBufferString(tt.req))
			resp, err := request.Login(c)
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
