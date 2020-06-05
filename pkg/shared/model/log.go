// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package sandpiper

import "github.com/labstack/echo/v4"

// Logger represents the logging interface implemented by each service
type Logger interface {
	// Log context, source, msg, error, params
	Log(echo.Context, string, string, error, map[string]interface{})
}
