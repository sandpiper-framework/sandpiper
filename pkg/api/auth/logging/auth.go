// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package auth

// auth service logger

import (
	"time"

	"github.com/labstack/echo/v4"

	"sandpiper/pkg/api/auth"
	"sandpiper/pkg/shared/model"
)

// ServiceLogger creates new logger wrapping the auth service
func ServiceLogger(svc auth.Service, logger sandpiper.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents auth logging service
type LogService struct {
	auth.Service
	logger sandpiper.Logger
}

const svc = "auth"

// Authenticate logging
func (ls *LogService) Authenticate(c echo.Context, user, password string) (resp *sandpiper.AuthToken, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			svc, "Authenticate request", err,
			map[string]interface{}{
				"req":  user,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Authenticate(c, user, password)
}

// Refresh logging
func (ls *LogService) Refresh(c echo.Context, req string) (resp *sandpiper.RefreshToken, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			svc, "Refresh request", err,
			map[string]interface{}{
				"req":  req,
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Refresh(c, req)
}

// Me logging
func (ls *LogService) Me(c echo.Context) (resp *sandpiper.User, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			svc, "Me request", err,
			map[string]interface{}{
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Me(c)
}

// Server logging
func (ls *LogService) Server(c echo.Context) (resp *sandpiper.Server) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			svc, "Server request", nil,
			map[string]interface{}{
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Server(c)
}
