package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

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

// Start starts echo server handling graceful shutdown, needs go1.8+.
func Start(e *echo.Echo, cfg *config.Server) {
	s := &http.Server{
		Addr:         cfg.Port,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Minute,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Minute,
	}
	e.Debug = cfg.Debug

	// Start server
	go func() {
		if err := e.StartServer(s); err != nil {
			e.Logger.Info("Shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
