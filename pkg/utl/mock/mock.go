package mock

import (
	"net/http/httptest"
	"time"

	"github.com/labstack/echo"
)

// TestTime is used for testing time fields
func TestTime(year int) time.Time {
	return time.Date(year, time.May, 19, 1, 2, 3, 4, time.UTC)
}

// HeaderValid is used for jwt testing
func HeaderValid() string {
	return "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidSI6ImpvaG5kb2UiLCJlIjoiam9obmRvZUBtYWlsLmNvbSIsInIiOjEsImMiOjEsImwiOjEsImV4cCI6NDEwOTMyMDg5NCwiaWF0IjoxNTE2MjM5MDIyfQ.FZYfTeKNulzfMM6IMSOHJKokgFACtL_e-MW7HOLLm1k"
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
