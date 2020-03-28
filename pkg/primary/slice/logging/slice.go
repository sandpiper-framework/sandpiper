// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package slice

// slice service logger

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/primary/slice"
	"autocare.org/sandpiper/pkg/shared/model"
)

// ServiceLogger creates new logger wrapping the slice service
func ServiceLogger(svc slice.Service, logger sandpiper.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents slice logging service
type LogService struct {
	slice.Service
	logger sandpiper.Logger
}

const source = "slice"

// Create logging
func (ls *LogService) Create(c echo.Context, req sandpiper.Slice) (resp *sandpiper.Slice, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "Create slice request", err,
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
func (ls *LogService) List(c echo.Context, tags *sandpiper.TagQuery, req *sandpiper.Pagination) (resp []sandpiper.Slice, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "List slice request", err,
			map[string]interface{}{
				"tags": tags,
				"req":  req,
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.List(c, tags, req)
}

// View logging
func (ls *LogService) View(c echo.Context, req uuid.UUID) (resp *sandpiper.Slice, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "View slice (by id) request", err,
			map[string]interface{}{
				"req":  req,
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.View(c, req)
}

// ViewByName logging
func (ls *LogService) ViewByName(c echo.Context, req string) (resp *sandpiper.Slice, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "View slice (by name) request", err,
			map[string]interface{}{
				"req":  req,
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.ViewByName(c, req)
}

// Delete logging
func (ls *LogService) Delete(c echo.Context, req uuid.UUID) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "Delete slice request", err,
			map[string]interface{}{
				"req":  req,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Delete(c, req)
}

// Update logging
func (ls *LogService) Update(c echo.Context, req *slice.Update) (resp *sandpiper.Slice, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "Update slice request", err,
			map[string]interface{}{
				"req":  req,
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Update(c, req)
}
