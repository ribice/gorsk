package jwt_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ribice/gorsk/pkg/utl/middleware/jwt"
	"github.com/ribice/gorsk/pkg/utl/model"

	"github.com/ribice/gorsk/pkg/utl/mock"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func echoHandler(mw ...echo.MiddlewareFunc) *echo.Echo {
	e := echo.New()
	for _, v := range mw {
		e.Use(v)
	}
	e.GET("/hello", hwHandler)
	return e
}

func hwHandler(c echo.Context) error {
	return c.String(200, "Hello World")
}

func TestMWFunc(t *testing.T) {
	cases := []struct {
		name       string
		wantStatus int
		header     string
		signMethod string
	}{
		{
			name:       "Empty header",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Header not containing Bearer",
			header:     "notBearer",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Invalid header",
			header:     mock.HeaderInvalid(),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Success",
			header:     mock.HeaderValid(),
			wantStatus: http.StatusOK,
		},
	}
	jwtMW := jwt.New("jwtsecret", "HS256", 60)
	ts := httptest.NewServer(echoHandler(jwtMW.MWFunc()))
	defer ts.Close()
	path := ts.URL + "/hello"
	client := &http.Client{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", path, nil)
			req.Header.Set("Authorization", tt.header)
			res, err := client.Do(req)
			if err != nil {
				t.Fatal("Cannot create http request")
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestGenerateToken(t *testing.T) {
	cases := []struct {
		name      string
		wantToken string
		algo      string
		req       *gorsk.User
	}{
		{
			name: "Invalid algo",
			algo: "invalid",
		},
		{
			name: "Success",
			algo: "HS256",
			req: &gorsk.User{
				Base: gorsk.Base{
					ID: 1,
				},
				Username: "johndoe",
				Email:    "johndoe@mail.com",
				Role: &gorsk.Role{
					AccessLevel: gorsk.SuperAdminRole,
				},
				CompanyID:  1,
				LocationID: 1,
			},
			wantToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.algo != "HS256" {
				assert.Panics(t, func() {
					jwt.New("jwtsecret", tt.algo, 60)
				}, "The code did not panic")
				return
			}
			jwt := jwt.New("jwtsecret", tt.algo, 60)
			str, _, err := jwt.GenerateToken(tt.req)
			assert.Nil(t, err)
			assert.Equal(t, tt.wantToken, strings.Split(str, ".")[0])
		})
	}
}
