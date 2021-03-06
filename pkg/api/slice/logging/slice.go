// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package slice

// slice service logger

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/api/slice"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/params"
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
func (ls *LogService) List(c echo.Context, p *params.Params, tags *params.TagQuery) (resp []sandpiper.Slice, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "List slice request", err,
			map[string]interface{}{
				"params": p,
				"tags":   tags,
				"resp":   resp,
				"took":   time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.List(c, p, tags)
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

// Refresh logging
func (ls *LogService) Refresh(c echo.Context, req uuid.UUID) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "Refresh slice request", err,
			map[string]interface{}{
				"req":  req,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Refresh(c, req)
}

// Lock logging
func (ls *LogService) Lock(c echo.Context, req uuid.UUID) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "Lock slice request", err,
			map[string]interface{}{
				"slice_id": req,
				"took":     time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Lock(c, req)
}

// Unlock logging
func (ls *LogService) Unlock(c echo.Context, req uuid.UUID) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "UnLock slice request", err,
			map[string]interface{}{
				"slice_id": req,
				"took":     time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Unlock(c, req)
}
