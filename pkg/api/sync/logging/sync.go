// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sync

// sync service logger

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"sandpiper/pkg/api/sync"
	"sandpiper/pkg/shared/model"
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
func (ls *LogService) Start(c echo.Context, req uuid.UUID) (err error) {
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

// Subscriptions logging
func (ls *LogService) Subscriptions(c echo.Context) (subs []sandpiper.Subscription, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "Sync Subscriptions request", err,
			map[string]interface{}{
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Subscriptions(c)
}

// Grains logging
func (ls *LogService) Grains(c echo.Context, sliceID uuid.UUID, briefFlag bool) (resp []sandpiper.Grain, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "Sync Grains request", err,
			map[string]interface{}{
				"slice-id": sliceID,
				"brief":    briefFlag,
				"resp":     fmt.Sprintf("Count: %d", len(resp)),
				"took":     time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Grains(c, sliceID, briefFlag)
}

// Process logging
func (ls *LogService) Process(c echo.Context) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "Process sync request", err,
			map[string]interface{}{
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Process(c)
}
