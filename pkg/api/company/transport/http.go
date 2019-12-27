// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package transport

// company routing functions

import (
	"github.com/labstack/echo/v4"
	"github.com/satori/go.uuid"
	"net/http"

	"autocare.org/sandpiper/internal/model"
	"autocare.org/sandpiper/pkg/api/company"
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
	sr.PATCH("/:id", h.update)
	sr.DELETE("/:id", h.delete)
}

// Custom errors
var (
	// ErrInvalidCompanyUUID indicates a malformed uuid
	ErrInvalidCompanyUUID = echo.NewHTTPError(http.StatusBadRequest, "invalid company uuid")
)

// Company create request
type createReq struct {
	Name   string `json:"name" validate:"required,min=3"`
	Active bool   `json:"active"`
}

func (h *HTTP) create(c echo.Context) error {
	r := new(createReq)
	if err := c.Bind(r); err != nil {
		return err
	}
	usr, err := h.svc.Create(c, sandpiper.Company{
		ID:     uuid.NewV4(),
		Name:   r.Name,
		Active: r.Active,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, usr)
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
	id, err := uuid.FromString(c.Param("id"))
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
	ID     uuid.UUID `json:"-"`
	Name   string    `json:"name,omitempty" validate:"omitempty,min=3"`
	Active bool      `json:"active,omitempty" validate:"omitempty"`
}

func (h *HTTP) update(c echo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return ErrInvalidCompanyUUID
	}
	req := new(updateReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	usr, err := h.svc.Update(c, &company.Update{
		ID:     id,
		Name:   req.Name,
		Active: req.Active,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, usr)
}

func (h *HTTP) delete(c echo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return ErrInvalidCompanyUUID
	}
	if err := h.svc.Delete(c, id); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}