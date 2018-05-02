package mock

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/ribice/gorsk/cmd/api/server"
)

// TestTime is used for testing time fields
func TestTime(year int) time.Time {
	return time.Date(year, time.May, 19, 1, 2, 3, 4, time.UTC)
}

// TestTimePtr is used for testing pointer time fields
func TestTimePtr(year int) *time.Time {
	t := time.Date(year, time.May, 19, 1, 2, 3, 4, time.UTC)
	return &t
}

// Str2Ptr converts string to pointer
func Str2Ptr(s string) *string {
	return &s
}

// EchoCtxWithKeys returns new Echo context with keys
func EchoCtxWithKeys(keys []string, values ...interface{}) echo.Context {
	e := echo.New()
	w := httptest.NewRecorder()
	c := e.NewContext(nil, w)
	for i, k := range keys {
		c.Set(k, values[i])
	}
	return c
}

// EchoCtx returns new Echo context, with validator and content type set
func EchoCtx(r *http.Request, w http.ResponseWriter) echo.Context {
	r.Header.Set("Content-Type", "application/json")
	e := echo.New()
	e.Validator = &server.CustomValidator{V: validator.New()}
	e.Binder = &server.CustomBinder{}
	return e.NewContext(r, w)
}

// HeaderValid is used for jwt testing
func HeaderValid() string {
	return "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidSI6ImpvaG5kb2UiLCJlIjoiam9obmRvZUBtYWlsLmNvbSIsInIiOjEsImMiOjEsImwiOjEsImV4cCI6NDEwOTMyMDg5NCwiaWF0IjoxNTE2MjM5MDIyfQ.8Fa8mhshx3tiQVzS5FoUXte5lHHC4cvaa_tzvcel38I"
}

// HeaderInvalid is used for jwt testing
func HeaderInvalid() string {
	return "Bearer eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidSI6ImpvaG5kb2UiLCJlIjoiam9obmRvZUBtYWlsLmNvbSIsInIiOjEsImMiOjEsImwiOjEsImV4cCI6NDEwOTMyMDg5NCwiaWF0IjoxNTE2MjM5MDIyfQ.7uPfVeZBkkyhICZSEINZfPo7ZsaY0NNeg0ebEGHuAvNjFvoKNn8dWYTKaZrqE1X4"
}
