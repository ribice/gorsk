package request_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ribice/gorsk/cmd/api/request"
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
