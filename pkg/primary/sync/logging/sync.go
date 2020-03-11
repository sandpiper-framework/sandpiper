// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sync

// sync service logger

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/primary/sync"
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

// Create logging
func (ls *LogService) Create(c echo.Context, req sandpiper.Sync) (resp *sandpiper.Sync, err error) {
	// todo: consider a "debug" level that shows entire req/resp
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "Create sync request", err,
			map[string]interface{}{
				"req":  req,
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Create(c, req)
}

// List logging
func (ls *LogService) List(c echo.Context, req *sandpiper.Pagination) (resp []sandpiper.Sync, err error) {
	// todo: consider a "debug" level that shows entire resp
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "List sync request", err,
			map[string]interface{}{
				"req":  req,
				"resp": fmt.Sprintf("Count: %d", len(resp)),
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.List(c, req)
}

// View logging
func (ls *LogService) View(c echo.Context, req int) (resp *sandpiper.Sync, err error) {
	// todo: consider a "debug" level that shows entire resp
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "View sync request", err,
			map[string]interface{}{
				"req":  req,
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.View(c, req)
}

// Delete logging
func (ls *LogService) Delete(c echo.Context, req int) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "Delete sync request", err,
			map[string]interface{}{
				"req":  req,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Delete(c, req)
}
