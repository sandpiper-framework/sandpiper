// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package transport

// sync routing functions

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/primary/sync"
	"autocare.org/sandpiper/pkg/shared/model"
)

// HTTP represents user http service
type HTTP struct {
	svc sync.Service
}

// NewHTTP creates new sync http service
func NewHTTP(svc sync.Service, er *echo.Group) {
	h := HTTP{svc}
	sr := er.Group("/sync")
	sr.POST("", h.create)
	sr.GET("", h.list)
	sr.GET("/:id", h.view)
	sr.DELETE("/:id", h.delete)
}

// Custom errors
var (
	// ErrInvalidID indicates a malformed uuid
	ErrInvalidID = echo.NewHTTPError(http.StatusBadRequest, "Invalid numeric sync id")
)

// Sync create request
type createReq struct {
	ID       int       `json:"id"` // optional
	SliceID  uuid.UUID `json:"slice_id" validate:"required"`
	Message  string    `json:"message" validate:"required"`
	Duration time.Time `json:"duration" validate:"required"`
}

func (h *HTTP) create(c echo.Context) error {
	r := new(createReq)

	if err := c.Bind(r); err != nil {
		return err
	}

	result, err := h.svc.Create(c, sandpiper.Sync{
		SliceID:  &r.SliceID,
		Message:  r.Message,
		Duration: r.Duration,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

type listResponse struct {
	Slices []sandpiper.Sync `json:"syncs"`
	Page   int              `json:"page"`
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

	return c.JSON(http.StatusOK, listResponse{result, p.Page})
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
