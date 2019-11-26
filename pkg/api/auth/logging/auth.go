package auth

import (
	"time"

	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/auth"
	"autocare.org/sandpiper/pkg/model"
)

// New creates new auth logging service
func New(svc auth.Service, logger sandpiper.Logger) *LogService {
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

const name = "auth"

// Authenticate logging
func (ls *LogService) Authenticate(c echo.Context, user, password string) (resp *sandpiper.AuthToken, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Authenticate request", err,
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
			name, "Refresh request", err,
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
			name, "Me request", err,
			map[string]interface{}{
				"resp": resp,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.Service.Me(c)
}
