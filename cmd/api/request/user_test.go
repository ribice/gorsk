package request_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ribice/gorsk/internal/mock"
	"github.com/stretchr/testify/assert"

	"github.com/ribice/gorsk/cmd/api/request"
)

func TestUserUpdate(t *testing.T) {
	cases := []struct {
		name     string
		id       string
		req      string
		wantErr  bool
		wantData *request.UpdateUser
	}{
		{
			name:    "Fail on ID param",
			wantErr: true,
			id:      "NaN",
			req:     `{}`,
		},
		{
			name:    "Fail on binding JSON",
			wantErr: true,
			id:      "1",
			req:     `{"first_name":"j","last_name":"okocha"}`,
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
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PATCH", "/", bytes.NewBufferString(tt.req))
			c := mock.EchoCtx(req, w)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			resp, err := request.UserUpdate(c)
			assert.Equal(t, tt.wantData, resp)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
