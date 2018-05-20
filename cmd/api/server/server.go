package server

import (
	"net/http"
	"time"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/middleware"
	"github.com/ribice/gorsk/cmd/api/config"
	"github.com/ribice/gorsk/cmd/api/mw"

	"github.com/labstack/echo"
)

// New instantates new Echo server
func New() *echo.Echo {
	e := echo.New()
	mw.Add(e, middleware.Logger(), middleware.Recover(),
		mw.CORS(), mw.SecureHeaders())
	e.GET("/", healthCheck)
	e.Validator = &CustomValidator{V: validator.New()}
	custErr := &customErrHandler{e: e}
	e.HTTPErrorHandler = custErr.handler
	e.Binder = &CustomBinder{}
	return e
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

// Start starts echo server
func Start(e *echo.Echo, cfg *config.Server) {
	e.Server.Addr = cfg.Port
	e.Server.ReadTimeout = time.Duration(cfg.ReadTimeout) * time.Minute
	e.Server.WriteTimeout = time.Duration(cfg.WriteTimeout) * time.Minute
	e.Debug = cfg.Debug
	e.Logger.Fatal(gracehttp.Serve(e.Server))
}
