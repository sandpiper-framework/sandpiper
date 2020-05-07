// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package transport

// sync routing functions

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"

	"sandpiper/pkg/api/sync"
)

// HTTP represents user http service
type HTTP struct {
	svc sync.Service
}

// NewHTTP creates new sync http service
func NewHTTP(svc sync.Service, er *echo.Group) {
	h := HTTP{svc}
	sr := er.Group("/sync")
	sr.POST("/:compid", h.start) // for secondary servers only
	sr.GET("", h.process)        // for primary servers only
	sr.GET("/subs", h.subs)
}

// Custom errors
var (
	// ErrInvalidURL indicates a malformed url
	ErrInvalidURL = echo.NewHTTPError(http.StatusBadRequest, "Invalid uuid")
)

func (h *HTTP) start(c echo.Context) error {
	id, err := uuid.Parse(c.Param("compid"))
	if err != nil {
		return ErrInvalidURL
	}
	if err := h.svc.Start(c, id); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (h *HTTP) process(c echo.Context) error {
	if err := h.svc.Process(c); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (h *HTTP) subs(c echo.Context) error {
	result, err := h.svc.Subscriptions(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}
