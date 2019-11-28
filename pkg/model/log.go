package sandpiper

import "github.com/labstack/echo/v4"

// Logger represents logging interface
type Logger interface {
	// context, source, msg, error, params
	Log(echo.Context, string, string, error, map[string]interface{})
}
