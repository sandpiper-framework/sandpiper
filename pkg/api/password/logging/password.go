// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package password

// password service logger

import (
	"time"

	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/password"
	"autocare.org/sandpiper/pkg/internal/model"
)

// ServiceLogger creates new logger wrapping the password service
func ServiceLogger(svc password.Service, logger sandpiper.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents password logging service
type LogService struct {
	password.Service
	logger sandpiper.Logger
}

const source = "password"

// Change password logging
func (ls *LogService) Change(c echo.Context, id int, oldPass, newPass string) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "Change password request", err,
			map[string]interface{}{
				"req":  id,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Change(c, id, oldPass, newPass)
}
