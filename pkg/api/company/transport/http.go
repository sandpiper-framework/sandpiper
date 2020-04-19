// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package transport

// company service routing functions

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/company"
	"autocare.org/sandpiper/pkg/shared/model"
)

// HTTP represents user http service
type HTTP struct {
	svc company.Service
}

// NewHTTP creates new company http service
func NewHTTP(svc company.Service, er *echo.Group) {
	h := HTTP{svc}
	sr := er.Group("/companies")
	sr.POST("", h.create)
	sr.GET("", h.list)
	sr.GET("/:id", h.view)
	sr.PUT("/:id", h.update) // not a PATCH, body must include *all* fields
	sr.DELETE("/:id", h.delete)
}

// Custom errors
var (
	// ErrInvalidCompanyUUID indicates a malformed uuid
	ErrInvalidCompanyUUID = echo.NewHTTPError(http.StatusBadRequest, "invalid company uuid")
)

// Company create request
type createReq struct {
	ID       uuid.UUID `json:"id"` // optional
	Name     string    `json:"name" validate:"required,min=3"`
	SyncAddr string    `json:"sync_addr"`
	Active   bool      `json:"active"`
}

func (r createReq) id() uuid.UUID {
	if r.ID == uuid.Nil {
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
	result, err := h.svc.Create(c, sandpiper.Company{
		ID:       r.id(),
		Name:     r.Name,
		SyncAddr: r.SyncAddr,
		Active:   r.Active,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

type listResponse struct {
	Companies []sandpiper.Company `json:"companies"`
	Page      int                 `json:"page"`
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
		return ErrInvalidCompanyUUID
	}
	result, err := h.svc.View(c, id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

// Company update request
type updateReq struct {
	ID       uuid.UUID `json:"-"`
	Name     string    `json:"name,omitempty" validate:"omitempty,min=3"`
	SyncAddr string    `json:"sync_addr"`
	Active   bool      `json:"active,omitempty" validate:"omitempty"`
}

func (h *HTTP) update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ErrInvalidCompanyUUID
	}
	req := new(updateReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	result, err := h.svc.Update(c, &company.Update{
		ID:       id,
		Name:     req.Name,
		SyncAddr: req.SyncAddr,
		Active:   req.Active,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func (h *HTTP) delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ErrInvalidCompanyUUID
	}
	if err := h.svc.Delete(c, id); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
