// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package company

// company service logger

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/google/uuid"

	"autocare.org/sandpiper/internal/model"
	"autocare.org/sandpiper/pkg/api/company"
)

// ServiceLogger creates new logger wrapping the company service
func ServiceLogger(svc company.Service, logger sandpiper.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents company logging service
type LogService struct {
	company.Service
	logger sandpiper.Logger
}

const source = "company"

// Create logging
func (ls *LogService) Create(c echo.Context, req sandpiper.Company) (resp *sandpiper.Company, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "Create company request", err,
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
func (ls *LogService) List(c echo.Context, req *sandpiper.Pagination) (resp []sandpiper.Company, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "List company request", err,
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
func (ls *LogService) View(c echo.Context, req uuid.UUID) (resp *sandpiper.Company, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "View company request", err,
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
			source, "Delete company request", err,
			map[string]interface{}{
				"req":  req,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Delete(c, req)
}

// Update logging
func (ls *LogService) Update(c echo.Context, req *company.Update) (resp *sandpiper.Company, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "Update company request", err,
			map[string]interface{}{
				"req":  req,
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Update(c, req)
}
