// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package setting

// setting service logger

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/api/setting"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
)

// ServiceLogger creates new logger wrapping the setting service
func ServiceLogger(svc setting.Service, logger sandpiper.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents setting logging service
type LogService struct {
	setting.Service
	logger sandpiper.Logger
}

const source = "setting"

// Create logging
func (ls *LogService) Create(c echo.Context, req *sandpiper.Setting) (resp *sandpiper.Setting, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "Create setting request", err,
			map[string]interface{}{
				"req":  req,
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Create(c, req)
}

// View logging
func (ls *LogService) View(c echo.Context) (resp *sandpiper.Setting, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "View setting request", err,
			map[string]interface{}{
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.View(c)
}

// Update logging
func (ls *LogService) Update(c echo.Context, req *setting.Update) (resp *sandpiper.Setting, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "Update subscription request", err,
			map[string]interface{}{
				"req":  req,
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Update(c, req)
}
