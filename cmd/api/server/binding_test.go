package server_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ribice/gorsk/cmd/api/server"
	"github.com/ribice/gorsk/internal/mock"
	"github.com/stretchr/testify/assert"
)

type Req struct {
	Name string `json:"name" validate:"required"`
}

func TestBind(t *testing.T) {
	cases := []struct {
		name     string
		req      string
		wantErr  bool
		wantData *Req
	}{
		{
			name:     "Fail on binding",
			wantErr:  true,
			req:      `"bleja"`,
			wantData: &Req{Name: ""},
		},
		{
			name:     "Fail on validation",
			wantErr:  true,
			wantData: &Req{Name: ""},
		},
		{
			name:     "Success",
			req:      `{"name":"John"}`,
			wantData: &Req{Name: "John"},
		},
	}
	b := &server.CustomBinder{}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "", bytes.NewBufferString(tt.req))
			c := mock.EchoCtx(req, w)
			r := new(Req)
			err := b.Bind(r, c)
			assert.Equal(t, tt.wantData, r)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}

}
