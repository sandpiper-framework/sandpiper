// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package transport

// slice routing functions

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/slice"
	"autocare.org/sandpiper/pkg/internal/model"
)

// HTTP represents user http service
type HTTP struct {
	svc slice.Service
}

// NewHTTP creates new slice http service
func NewHTTP(svc slice.Service, er *echo.Group) {
	h := HTTP{svc}
	sr := er.Group("/slices")
	sr.POST("", h.create)
	sr.GET("", h.list)
	sr.GET("/:id", h.view)
	sr.PATCH("/:id", h.update)
	sr.DELETE("/:id", h.delete)
}

// Custom errors
var (
	// ErrInvalidSliceUUID indicates a malformed uuid
	ErrInvalidSliceUUID = echo.NewHTTPError(http.StatusBadRequest, "invalid slice uuid")
)

// Slice create request
type createReq struct {
	ID           uuid.UUID         `json:"id"` // optional
	Name         string            `json:"name" validate:"required,min=3"`
	ContentHash  string            `json:"content_hash"`
	ContentCount uint              `json:"content_count"`
	ContentDate  time.Time         `json:"content_date"`
	Metadata     sandpiper.MetaMap `json:"metadata"`
}

func (r createReq) id() uuid.UUID {
	var nilUUID = uuid.UUID{}
	if r.ID == nilUUID {
		return uuid.New()
	}
	return r.ID
}

// create populates createReq from body json adding UUID if not provided
func (h *HTTP) create(c echo.Context) error {
	r := new(createReq)

	if err := c.Bind(r); err != nil {
		return err
	}

	result, err := h.svc.Create(c, sandpiper.Slice{
		ID:           r.id(),
		Name:         r.Name,
		ContentHash:  r.ContentHash,
		ContentCount: r.ContentCount,
		ContentDate:  r.ContentDate,
		Metadata:     r.Metadata,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

type listResponse struct {
	Slices []sandpiper.Slice `json:"slices"`
	Page   int               `json:"page"`
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
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ErrInvalidSliceUUID
	}

	result, err := h.svc.View(c, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

// Slice update request
type updateReq struct {
	ID           uuid.UUID `json:"-"`
	Name         string    `json:"name,omitempty" validate:"omitempty,min=3"`
	ContentHash  string    `json:"content_hash,omitempty" validate:"omitempty,min=2"`
	ContentCount uint      `json:"content_count,omitempty"`
	ContentDate  time.Time `json:"content_date,omitempty"`
}

func (h *HTTP) update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ErrInvalidSliceUUID
	}

	req := new(updateReq)
	if err := c.Bind(req); err != nil {
		return err
	}

	result, err := h.svc.Update(c, &slice.Update{
		ID:           id,
		Name:         req.Name,
		ContentHash:  req.ContentHash,
		ContentCount: req.ContentCount,
		ContentDate:  req.ContentDate,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func (h *HTTP) delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ErrInvalidSliceUUID
	}

	if err := h.svc.Delete(c, id); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
