// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package zlog

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// Log represents zerolog logger
type Log struct {
	logger *zerolog.Logger
}

// New instantiates new zero logger with "Info" Level default messaging.
// Set log level to "warn" (and above) if `logInfoLevel` is false.
func New(logInfoLevel bool) *Log {
	var lvl zerolog.Level = zerolog.InfoLevel // default

	if !logInfoLevel {
		lvl = zerolog.WarnLevel
	}

	z := zerolog.New(os.Stdout).Level(lvl)

	return &Log{
		logger: &z,
	}
}

// Log implements the sandpiper.Logger interface (using zerolog)
func (z *Log) Log(ctx echo.Context, source, msg string, err error, params map[string]interface{}) {

	if params == nil {
		params = make(map[string]interface{})
	}

	params["service"] = source

	if id, ok := ctx.Get("id").(int); ok {
		params["user_id"] = id
		params["user"] = ctx.Get("username").(string)
	}

	if err != nil {
		params["error"] = err
		z.logger.Error().Fields(params).Msg(msg)
		return
	}

	z.logger.Info().Fields(params).Msg(msg)
}
