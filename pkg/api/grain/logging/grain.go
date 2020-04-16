// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package grain

// grain service logger

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/grain"
	"autocare.org/sandpiper/pkg/shared/model"
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
func (ls *LogService) Create(c echo.Context, replaceFlag bool, req *sandpiper.Grain) (resp *sandpiper.Grain, err error) {
	// todo: consider a "debug" level that shows entire req/resp
	defer func(begin time.Time) {
		var g *sandpiper.Grain
		// suppress payload in log for req and resp
		req.Payload = sandpiper.PayloadNil
		if resp != nil {
			g = &sandpiper.Grain{
				ID:       resp.ID,
				SliceID:  resp.SliceID,
				Key:      resp.Key,
				Source:   resp.Source,
				Encoding: resp.Encoding,
			}
		}
		ls.logger.Log(
			c,
			source, "Create grain request", err,
			map[string]interface{}{
				"req":     req,
				"replace": replaceFlag,
				"resp":    g,
				"took":    time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Create(c, replaceFlag, req)
}

// List logging
func (ls *LogService) List(c echo.Context, payload bool, req *sandpiper.Pagination) (resp []sandpiper.Grain, err error) {
	// todo: consider a "debug" level that shows entire resp
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "List grain request", err,
			map[string]interface{}{
				"req":  req,
				"resp": fmt.Sprintf("Count: %d", len(resp)),
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.List(c, payload, req)
}

// View logging
func (ls *LogService) View(c echo.Context, req uuid.UUID) (resp *sandpiper.Grain, err error) {
	// todo: consider a "debug" level that shows entire resp (including payload)
	defer func(begin time.Time) {
		var g *sandpiper.Grain
		if resp != nil {
			// suppress payload in log
			g = &sandpiper.Grain{
				ID:       resp.ID,
				SliceID:  resp.SliceID,
				Key:      resp.Key,
				Source:   resp.Source,
				Encoding: resp.Encoding,
			}
		}
		ls.logger.Log(
			c,
			source, "View grain request", err,
			map[string]interface{}{
				"req":  req,
				"resp": g,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.View(c, req)
}

// ViewByKeys logging
func (ls *LogService) ViewByKeys(c echo.Context, sliceID uuid.UUID, grainKey string, payloadFlag bool) (resp *sandpiper.Grain, err error) {
	defer func(begin time.Time) {
		var g *sandpiper.Grain
		if resp != nil {
			g = &sandpiper.Grain{
				ID:      resp.ID,
				SliceID: resp.SliceID,
				Key:     resp.Key,
				Source:  resp.Source,
			}
		}
		ls.logger.Log(
			c,
			source, "Exists grain request", err,
			map[string]interface{}{
				"slice_id":  sliceID,
				"grain_key": grainKey,
				"resp":      g,
				"took":      time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.ViewByKeys(c, sliceID, grainKey, payloadFlag)
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
