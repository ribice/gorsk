package request_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ribice/gorsk/cmd/api/request"
	"github.com/ribice/gorsk/internal/mock"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	cases := []struct {
		name     string
		req      string
		wantErr  bool
		wantData *request.Credentials
	}{
		{
			name:    "Fail on binding JSON",
			wantErr: true,
			req:     `{"username":"juzernejm"}`,
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

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "", bytes.NewBufferString(tt.req))
			c := mock.EchoCtx(req, w)
			resp, err := request.Login(c)
			assert.Equal(t, tt.wantData, resp)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
