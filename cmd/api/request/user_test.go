package request_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ribice/gorsk/internal/mock"
	"github.com/stretchr/testify/assert"

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
				assert.Equal(t, tt.e.wantStatus, w.Code)
				assert.Equal(t, tt.e.wantResp, w.Body.String())
			}
			assert.Equal(t, tt.wantData, resp)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantErr, c.IsAborted())
		})
	}
}
