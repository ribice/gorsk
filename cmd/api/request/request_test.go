package request_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/ribice/gorsk/cmd/api/request"
	"github.com/ribice/gorsk/internal/mock"
	"github.com/stretchr/testify/assert"
)

func TestPaginate(t *testing.T) {
	cases := []struct {
		name     string
		req      string
		wantErr  bool
		wantData *request.Pagination
	}{
		{
			name:    "Fail on binding JSON",
			wantErr: true,
			req:     `/?limit=50&page=-1`,
		},
		{
			name: "Test default limit",
			req:  `/?limit=0`,
			wantData: &request.Pagination{
				Limit: 100,
			},
		},
		{
			name: "Test max limit",
			req:  `/?limit=2222&page=2`,
			wantData: &request.Pagination{
				Limit:  1000,
				Offset: 2000,
				Page:   2,
			},
		},
		{
			name: "Test default",
			req:  `/?limit=200&page=2`,
			wantData: &request.Pagination{
				Limit:  200,
				Offset: 400,
				Page:   2,
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest("GET", tt.req, nil)
			if err != nil {
				t.Error("Could not create http request")
			}
			c := mock.EchoCtx(req, w)
			resp, err := request.Paginate(c)
			assert.Equal(t, tt.wantData, resp)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
func TestID(t *testing.T) {
	cases := []struct {
		name     string
		id       string
		wantErr  bool
		wantData int
	}{
		{
			name:    "EmptyID",
			wantErr: true,
		},
		{
			name:    "ID Not a Number",
			id:      "NaN",
			wantErr: true,
		},
		{
			name:     "Success",
			wantData: 1,
			id:       "1",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(echo.GET, "/", nil)
			c := mock.EchoCtx(req, w)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			resp, err := request.ID(c)
			assert.Equal(t, tt.wantData, resp)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
