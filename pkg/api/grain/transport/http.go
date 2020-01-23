// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package transport

// grain routing functions

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/grain"
	"autocare.org/sandpiper/pkg/shared/model"
)

// HTTP represents user http service
type HTTP struct {
	svc grain.Service
}

// NewHTTP creates new grain http service
func NewHTTP(svc grain.Service, er *echo.Group) {
	h := HTTP{svc}
	sr := er.Group("/grains")
	sr.POST("", h.create)
	sr.GET("", h.list)
	sr.GET("/:id", h.view)
	sr.DELETE("/:id", h.delete)
}

// Custom errors
var (
	// ErrInvalidSliceUUID indicates a malformed uuid
	ErrInvalidSliceUUID = echo.NewHTTPError(http.StatusBadRequest, "Invalid grain uuid")
)

// Grain create request
type createReq struct {
	ID       uuid.UUID `json:"id"` // optional
	SliceID  uuid.UUID `json:"slice_id" validate:"required"`
	Type     string    `json:"grain_type" validate:"required"`
	Key      string    `json:"grain_key" validate:"required"`
	Encoding string    `json:"encoding" validate:"required"`
	Payload  []byte    `json:"payload" validate:"required"`
}

func (r createReq) id() uuid.UUID {
	var nilUUID = uuid.UUID{}
	if r.ID == nilUUID {
		return uuid.New()
	}
	return r.ID
}

func (h *HTTP) create(c echo.Context) error {
	r := new(createReq)

	if err := c.Bind(r); err != nil {
		return err
	}

	result, err := h.svc.Create(c, sandpiper.Grain{
		ID:       r.id(),
		SliceID:  &r.SliceID,
		Type:     r.Type,
		Key:      r.Key,
		Encoding: r.Encoding,
		Payload:  r.Payload,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

type listResponse struct {
	Slices []sandpiper.Grain `json:"grains"`
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
