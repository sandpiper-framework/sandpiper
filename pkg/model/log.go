package sandpiper

import "github.com/labstack/echo/v4"

// Logger represents the logging interface implemented by each service
type Logger interface {
	// context, source, msg, error, params
	Log(echo.Context, string, string, error, map[string]interface{})
}
