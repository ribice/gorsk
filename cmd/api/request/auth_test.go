package request_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ribice/gorsk/cmd/api/request"
	"github.com/stretchr/testify/assert"
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
				assert.Equal(t, tt.e.wantStatus, w.Code)
				assert.Equal(t, tt.e.wantResp, w.Body.String())
			}
			assert.Equal(t, tt.wantData, resp)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantErr, c.IsAborted())
		})
	}
}
