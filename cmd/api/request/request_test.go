package request_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ribice/gorsk/cmd/api/request"
	"github.com/stretchr/testify/assert"
)

func TestPaginate(t *testing.T) {
	type errResp struct {
		wantStatus int
		wantResp   string
	}
	cases := []struct {
		name     string
		e        *errResp
		req      string
		wantErr  bool
		wantData *request.Pagination
	}{
		{
			name:    "Fail on binding JSON",
			wantErr: true,
			req:     `?limit=50&page=-1`,
			e: &errResp{
				wantStatus: http.StatusBadRequest,
				wantResp:   `{"message":["Page's value or length is less than allowed"]}`,
			},
		},
		{
			name: "Test default limit",
			req:  `?limit=0`,
			wantData: &request.Pagination{
				Limit: 100,
			},
		},
		{
			name: "Test max limit",
			req:  `?limit=2222&page=2`,
			wantData: &request.Pagination{
				Limit:  1000,
				Offset: 2000,
				Page:   2,
			},
		},
		{
			name: "Test default",
			req:  `?limit=200&page=2`,
			wantData: &request.Pagination{
				Limit:  200,
				Offset: 400,
				Page:   2,
			},
		},
	}
	gin.SetMode(gin.TestMode)
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("POST", tt.req, nil)
			resp, err := request.Paginate(c)
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
func TestID(t *testing.T) {
	type errResp struct {
		wantStatus int
		wantResp   string
	}
	cases := []struct {
		name     string
		e        *errResp
		id       string
		wantErr  bool
		wantData int
	}{
		{
			name:    "EmptyID",
			wantErr: true,
			e: &errResp{
				wantStatus: http.StatusBadRequest,
			},
		},
		{
			name:     "ID Not a Number",
			id:       "NaN",
			wantErr:  true,
			wantData: 0,
			e: &errResp{
				wantStatus: http.StatusBadRequest,
			},
		},
		{
			name:     "Success",
			wantData: 1,
			id:       "1",
		},
	}
	gin.SetMode(gin.TestMode)
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{gin.Param{Key: "id", Value: tt.id}}
			resp, err := request.ID(c)
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
