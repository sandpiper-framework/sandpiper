// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package activity

// activity service logger

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/api/activity"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/params"
)

// ServiceLogger creates new logger wrapping the activity service
func ServiceLogger(svc activity.Service, logger sandpiper.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents activity logging service
type LogService struct {
	activity.Service
	logger sandpiper.Logger
}

const source = "activity"

// Create logging
func (ls *LogService) Create(c echo.Context, req sandpiper.Activity) (resp *sandpiper.Activity, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "Create activity request", err,
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
func (ls *LogService) List(c echo.Context, req *params.Params) (resp []sandpiper.Activity, err error) {
	// todo: consider a "debug" level that shows entire resp
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "List activity request", err,
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
func (ls *LogService) View(c echo.Context, req int) (resp *sandpiper.Activity, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "View activity request", err,
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
			source, "Delete activity request", err,
			map[string]interface{}{
				"req":  req,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Delete(c, req)
}
