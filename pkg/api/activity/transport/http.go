// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package transport

// activity routing

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"sandpiper/pkg/api/activity"
	"sandpiper/pkg/shared/model"
)

// HTTP represents user http service
type HTTP struct {
	svc activity.Service
}

// NewHTTP creates new activity http service
func NewHTTP(svc activity.Service, er *echo.Group) {
	h := HTTP{svc}
	sr := er.Group("/activity")
	sr.POST("", h.create)
	sr.GET("", h.list)
	sr.GET("/:id", h.view)
	sr.DELETE("/:id", h.delete)
}

// Custom errors
var (
	// ErrInvalidID indicates a malformed uuid
	ErrInvalidID = echo.NewHTTPError(http.StatusBadRequest, "Invalid numeric activity id")
)

// activity create request
type createReq struct {
	CompanyID uuid.UUID `json:"company_id" validate:"required"`
	SliceID   uuid.UUID `json:"slice_id" validate:"required"`
	Success   bool      `json:"success" validate:"required"`
	Message   string    `json:"message" validate:"required"`
	Duration  time.Time `json:"duration" validate:"required"`
}

func (h *HTTP) create(c echo.Context) error {
	r := new(createReq)

	if err := c.Bind(r); err != nil {
		return err
	}

	result, err := h.svc.Create(c, sandpiper.Activity{
		CompanyID: r.CompanyID,
		SliceID:   r.SliceID,
		Success:   r.Success,
		Message:   r.Message,
		Duration:  r.Duration,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, result)
}

func (h *HTTP) list(c echo.Context) error {
	p := new(sandpiper.PaginationReq)
	if err := c.Bind(p); err != nil {
		return err
	}

	result, err := h.svc.List(c, p.Transform())

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, sandpiper.ActivityPaginated{Syncs: result, Page: p.Page})
}

func (h *HTTP) view(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ErrInvalidID
	}

	result, err := h.svc.View(c, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func (h *HTTP) delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ErrInvalidID
	}

	if err := h.svc.Delete(c, id); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
