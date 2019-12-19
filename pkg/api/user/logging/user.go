package user

// user service logger

import (
	"time"

	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/internal/model"
	"autocare.org/sandpiper/pkg/api/user"
)

// New creates new user logging service
func New(svc user.Service, logger sandpiper.Logger) *LogService {
	return &LogService{
		Service: svc,
		logger:  logger,
	}
}

// LogService represents user logging service
type LogService struct {
	user.Service
	logger sandpiper.Logger
}

const source = "user"

// Create logging
func (ls *LogService) Create(c echo.Context, req sandpiper.User) (resp *sandpiper.User, err error) {
	defer func(begin time.Time) {
		req.Password = "xxx-redacted-xxx"
		ls.logger.Log(
			c,
			source, "Create user request", err,
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
func (ls *LogService) List(c echo.Context, req *sandpiper.Pagination) (resp []sandpiper.User, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "List user request", err,
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
func (ls *LogService) View(c echo.Context, req int) (resp *sandpiper.User, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "View user request", err,
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
			source, "Delete user request", err,
			map[string]interface{}{
				"req":  req,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Delete(c, req)
}

// Update logging
func (ls *LogService) Update(c echo.Context, req *user.Update) (resp *sandpiper.User, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			source, "Update user request", err,
			map[string]interface{}{
				"req":  req,
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Update(c, req)
}
