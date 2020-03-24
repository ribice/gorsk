package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"github.com/ribice/gorsk"
	"github.com/ribice/gorsk/pkg/utl/middleware/auth"
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

type tokenParser struct {
	ParseTokenFn func(string) (*jwt.Token, error)
}

func (t tokenParser) ParseToken(s string) (*jwt.Token, error) {
	if s == "" {
		return nil, gorsk.ErrGeneric
	}
	return &jwt.Token{
		Raw:    "abcd",
		Method: jwt.SigningMethodHS256,
		Claims: jwt.MapClaims{
			"c":   1.0,
			"e":   "johndoe@mail.com",
			"exp": 1581773411,
			"id":  1.0,
			"l":   1.0,
			"r":   100.0,
			"u":   "johndoe",
		},
		Valid: true,
	}, nil
}

func TestMWFunc(t *testing.T) {
	cases := map[string]struct {
		wantStatus int
		header     string
		signMethod string
	}{
		"Empty header": {
			wantStatus: http.StatusUnauthorized,
		},
		"Success": {
			header:     "Bearer 123",
			wantStatus: http.StatusOK,
		},
	}
	ts := httptest.NewServer(echoHandler(auth.Middleware(tokenParser{})))
	defer ts.Close()
	path := ts.URL + "/hello"
	client := &http.Client{}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
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
