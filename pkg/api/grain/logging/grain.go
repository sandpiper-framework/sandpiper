// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package grain

// grain service logger

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/internal/model"
	"autocare.org/sandpiper/pkg/api/grain"
)

// ServiceLogger creates new logger wrapping the grain service
func ServiceLogger(svc grain.Service, logger sandpiper.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents grain logging service
type LogService struct {
	grain.Service
	logger sandpiper.Logger
}

const source = "grain"

// Create logging
func (ls *LogService) Create(c echo.Context, req sandpiper.Grain) (resp *sandpiper.Grain, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "Create grain request", err,
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
func (ls *LogService) List(c echo.Context, req *sandpiper.Pagination) (resp []sandpiper.Grain, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "List grain request", err,
			map[string]interface{}{
				"req":  req,
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.List(c, req)
}

// View logging
func (ls *LogService) View(c echo.Context, req uuid.UUID) (resp *sandpiper.Grain, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "View grain request", err,
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
func (ls *LogService) Delete(c echo.Context, req uuid.UUID) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "Delete grain request", err,
			map[string]interface{}{
				"req":  req,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Delete(c, req)
}
