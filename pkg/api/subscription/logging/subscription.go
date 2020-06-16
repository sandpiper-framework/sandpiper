// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package subscription

// subscription service logger

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/api/subscription"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/params"
)

// ServiceLogger creates new logger wrapping the subscription service
func ServiceLogger(svc subscription.Service, logger sandpiper.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents subscription logging service
type LogService struct {
	subscription.Service
	logger sandpiper.Logger
}

const source = "subscription"

// Create logging
func (ls *LogService) Create(c echo.Context, req sandpiper.Subscription) (resp *sandpiper.Subscription, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "Create subscription request", err,
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
func (ls *LogService) List(c echo.Context, req *params.Params) (resp []sandpiper.Subscription, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "List subscription request", err,
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
func (ls *LogService) View(c echo.Context, req sandpiper.Subscription) (resp *sandpiper.Subscription, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "View subscription request", err,
			map[string]interface{}{
				"req-id":         req.SubID,
				"req-slice_id":   req.SliceID,
				"req-company_id": req.CompanyID,
				"req-name":       req.Name,
				"resp":           resp,
				"took":           time.Since(begin),
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
			source, "Delete subscription request", err,
			map[string]interface{}{
				"req":  req,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Delete(c, req)
}

// Update logging
func (ls *LogService) Update(c echo.Context, req *subscription.Update) (resp *sandpiper.Subscription, err error) {
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
