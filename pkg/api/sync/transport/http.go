// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package transport

// sync routing functions

// Some functionality intentionally duplicates other services because we plan on
// changing this service to use websockets eventually and it will be easier if
// isolated now. We also don't want pagination of these resources.

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/api/sync"
)

// Custom errors
var (
	// ErrInvalidSliceUUID indicates an improperly formed UUID
	ErrInvalidSliceUUID = echo.NewHTTPError(http.StatusBadRequest, "Invalid slice uuid")
)

// HTTP represents user http service
type HTTP struct {
	svc sync.Service
}

// NewHTTP creates new sync http service
func NewHTTP(svc sync.Service, er *echo.Group) {
	h := HTTP{svc}
	sr := er.Group("/sync")
	sr.POST("/:compid", h.start)   // for secondary servers only
	sr.GET("", h.process)          // for primary servers only
	sr.GET("/subs", h.subs)        // get my subscriptions
	sr.GET("/slice/:id", h.grains) // ?brief=yes|no
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

func (h *HTTP) grains(c echo.Context) error {
	var briefFlag bool

	sliceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ErrInvalidSliceUUID
	}
	if c.QueryParam("brief") == "yes" {
		briefFlag = true
	}
	result, err := h.svc.Grains(c, sliceID, briefFlag)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}
