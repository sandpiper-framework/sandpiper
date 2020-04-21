// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sync

// sync service logger

import (
	"net/url"
	"time"

	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/sync"
	"autocare.org/sandpiper/pkg/shared/model"
)

// ServiceLogger creates new logger wrapping the sync service
func ServiceLogger(svc sync.Service, logger sandpiper.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents sync logging service
type LogService struct {
	sync.Service
	logger sandpiper.Logger
}

const source = "sync"

// Start logging
func (ls *LogService) Start(c echo.Context, req *url.URL) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "Start sync request", err,
			map[string]interface{}{
				"req":  req,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Start(c, req)
}

// Connect logging
func (ls *LogService) Connect(c echo.Context) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "Connect sync request", err,
			map[string]interface{}{
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Connect(c)
}
