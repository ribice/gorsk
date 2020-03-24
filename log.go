package gorsk

import "github.com/labstack/echo"

// Logger represents logging interface
type Logger interface {
	// source, msg, error, params
	Log(echo.Context, string, string, error, map[string]interface{})
}
