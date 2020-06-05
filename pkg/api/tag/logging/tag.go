// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package tag

// tag service logger

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"sandpiper/pkg/api/tag"
	"sandpiper/pkg/shared/model"
)

// ServiceLogger creates new logger wrapping the subscription service
func ServiceLogger(svc tag.Service, logger sandpiper.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents subscription logging service
type LogService struct {
	tag.Service
	logger sandpiper.Logger
}

const source = "tag"

// Create logging
func (ls *LogService) Create(c echo.Context, req sandpiper.Tag) (resp *sandpiper.Tag, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "Create tag request", err,
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
func (ls *LogService) List(c echo.Context, req *sandpiper.Pagination) (resp []sandpiper.Tag, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "List tag request", err,
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
func (ls *LogService) View(c echo.Context, req int) (resp *sandpiper.Tag, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "View tag request", err,
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
			source, "Delete tag request", err,
			map[string]interface{}{
				"req":  req,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Delete(c, req)
}

// Update logging
func (ls *LogService) Update(c echo.Context, req *tag.Update) (resp *sandpiper.Tag, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "Update tag request", err,
			map[string]interface{}{
				"req":  req,
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Update(c, req)
}

// Assign logging
func (ls *LogService) Assign(c echo.Context, tagID int, sliceID uuid.UUID) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "assign request", err,
			map[string]interface{}{
				"tagID":   tagID,
				"sliceID": sliceID,
				"took":    time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Assign(c, tagID, sliceID)
}

// Remove logging
func (ls *LogService) Remove(c echo.Context, tagID int, sliceID uuid.UUID) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "assign request", err,
			map[string]interface{}{
				"tagID":   tagID,
				"sliceID": sliceID,
				"took":    time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Remove(c, tagID, sliceID)
}
