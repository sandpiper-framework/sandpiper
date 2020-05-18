// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package transport

// grain routing

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"sandpiper/pkg/api/grain"
	"sandpiper/pkg/shared/model"
	"sandpiper/pkg/shared/payload"
)

// HTTP represents user http service
type HTTP struct {
	svc grain.Service
}

// NewHTTP creates new grain http service
func NewHTTP(svc grain.Service, er *echo.Group) {
	h := HTTP{svc}
	sr := er.Group("/grains")
	sr.POST("", h.create) // ?replace=[yes/no*]
	sr.GET("", h.list)    // ?payload=[yes/no*]
	sr.GET("/slice/:id", h.listBySlice)
	sr.GET("/:id", h.view)
	sr.GET("/:sliceid/:grainkey", h.viewByKeys) // ?payload=[yes/no*]
	sr.DELETE("/:id", h.delete)
}

// Custom errors
var (
	ErrInvalidGrainUUID = echo.NewHTTPError(http.StatusBadRequest, "Invalid grain uuid")
	ErrInvalidSliceUUID = echo.NewHTTPError(http.StatusBadRequest, "Invalid slice uuid")
)

// Grain create request
type createReq struct {
	ID       uuid.UUID           `json:"id"` // optional
	SliceID  uuid.UUID           `json:"slice_id" validate:"required"`
	Key      string              `json:"grain_key" validate:"required"`
	Source   string              `json:"source"`
	Encoding string              `json:"encoding" validate:"required"`
	Payload  payload.PayloadData `json:"payload" validate:"required"`
}

func (r createReq) id() uuid.UUID {
	if r.ID == uuid.Nil {
		return uuid.New()
	}
	return r.ID
}

func (h *HTTP) create(c echo.Context) error {
	var replaceFlag = false

	if c.QueryParam("replace") == "yes" {
		replaceFlag = true
	}

	r := new(createReq)
	if err := c.Bind(r); err != nil {
		return err
	}

	result, err := h.svc.Create(c, replaceFlag, &sandpiper.Grain{
		ID:       r.id(),
		SliceID:  &r.SliceID,
		Key:      r.Key,
		Source:   r.Source,
		Encoding: r.Encoding,
		Payload:  r.Payload,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, result)
}

func (h *HTTP) list(c echo.Context) error {
	var includePayload = false

	if c.QueryParam("payload") == "yes" {
		includePayload = true
	}

	p := new(sandpiper.PaginationReq)
	if err := c.Bind(p); err != nil {
		return err
	}

	result, err := h.svc.List(c, includePayload, p.Transform())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, sandpiper.GrainsPaginated{Grains: result, Page: p.Page})
}

func (h *HTTP) listBySlice(c echo.Context) error {
	var includePayload = false

	sliceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ErrInvalidSliceUUID
	}

	if c.QueryParam("payload") == "yes" {
		includePayload = true
	}

	p := new(sandpiper.PaginationReq)
	if err := c.Bind(p); err != nil {
		return err
	}

	result, err := h.svc.ListBySlice(c, sliceID, includePayload, p.Transform())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, sandpiper.GrainsPaginated{Grains: result, Page: p.Page})
}

func (h *HTTP) view(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ErrInvalidGrainUUID
	}

	result, err := h.svc.View(c, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func (h *HTTP) viewByKeys(c echo.Context) error {
	var includePayload = false

	SliceID, err := uuid.Parse(c.Param("sliceid"))
	if err != nil {
		return ErrInvalidSliceUUID
	}

	if c.QueryParam("payload") == "yes" {
		includePayload = true
	}

	result, err := h.svc.ViewByKeys(c, SliceID, c.Param("grainkey"), includePayload)
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
