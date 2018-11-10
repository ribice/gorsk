package zlog

import (
	"os"

	"github.com/labstack/echo"

	"github.com/rs/zerolog"
)

// Log represents zerolog logger
type Log struct {
	logger *zerolog.Logger
}

// New instantiates new zero logger
func New() *Log {
	z := zerolog.New(os.Stdout)
	return &Log{
		logger: &z,
	}
}

// Log logs using zerolog
func (z *Log) Log(ctx echo.Context, source, msg string, err error, params map[string]interface{}) {

	if params == nil {
		params = make(map[string]interface{})
	}

	params["source"] = source

	if id, ok := ctx.Get("id").(int); ok {
		params["id"] = id
		params["user"] = ctx.Get("username").(string)
	}

	if err != nil {
		params["error"] = err
		z.logger.Error().Fields(params).Msg(msg)
		return
	}

	z.logger.Info().Fields(params).Msg(msg)
}
