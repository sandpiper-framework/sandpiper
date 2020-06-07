// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package transport

// subscription service routing functions

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/api/subscription"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
)

// HTTP represents user http service
type HTTP struct {
	svc subscription.Service
}

// NewHTTP creates new subscription http service
func NewHTTP(svc subscription.Service, er *echo.Group) {
	h := HTTP{svc}

	er.GET("/companies/:id/subs", h.listByCompany)
	er.GET("/companies/:id/subs/:sliceid", h.viewByCompany)

	er.POST("/subs", h.create)
	er.GET("/subs", h.list)
	er.GET("/subs/:id", h.view)
	er.GET("/subs/name/:name", h.viewByName)
	er.PUT("/subs/:id", h.update) // not a PATCH, body must include *all* fields
	er.DELETE("/subs/:id", h.delete)
}

// Custom errors
var (
	ErrInvalidCompanyUUID      = echo.NewHTTPError(http.StatusBadRequest, "malformed company uuid")
	ErrInvalidSliceUUID        = echo.NewHTTPError(http.StatusBadRequest, "malformed slice uuid")
	ErrInvalidSubscriptionUUID = echo.NewHTTPError(http.StatusBadRequest, "invalid subscription uuid")
	ErrMissingSubscriptionName = echo.NewHTTPError(http.StatusBadRequest, "missing required subscription name")
)

// Subscription create request
type createReq struct {
	ID          uuid.UUID `json:"id"` // optional
	SliceID     uuid.UUID `json:"slice_id" validate:"required"`
	CompanyID   uuid.UUID `json:"company_id" validate:"required"`
	Name        string    `json:"name" validate:"required,min=3"`
	Description string    `json:"description"`
	Active      bool      `json:"active"`
}

func (r createReq) id() uuid.UUID {
	if r.ID == uuid.Nil {
		return uuid.New()
	}
	return r.ID
}

// create populates createReq from supplied json body
func (h *HTTP) create(c echo.Context) error {
	r := new(createReq)
	if err := c.Bind(r); err != nil {
		return err
	}
	result, err := h.svc.Create(c, sandpiper.Subscription{
		SubID:       r.id(),
		SliceID:     r.SliceID,
		CompanyID:   r.CompanyID,
		Name:        r.Name,
		Description: r.Description,
		Active:      r.Active,
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
	return c.JSON(http.StatusOK, sandpiper.SubsPaginated{Subs: result, Page: p.Page})
}

func (h *HTTP) listByCompany(c echo.Context) error {
	// todo: fix this... pass in companyID as a filter
	//id, err := uuid.Parse(c.Param("id"))
	//if err != nil {
	//	return ErrInvalidCompanyUUID
	//}
	p := new(sandpiper.PaginationReq)
	if err := c.Bind(p); err != nil {
		return err
	}
	result, err := h.svc.List(c, p.Transform())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, sandpiper.SubsPaginated{Subs: result, Page: p.Page})
}

func (h *HTTP) view(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ErrInvalidSubscriptionUUID
	}

	sub := sandpiper.Subscription{SubID: id}
	result, err := h.svc.View(c, sub)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func (h *HTTP) viewByCompany(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ErrInvalidCompanyUUID
	}
	sliceID, err := uuid.Parse(c.Param("sliceid"))
	if err != nil {
		return ErrInvalidSliceUUID
	}
	sub := sandpiper.Subscription{CompanyID: id, SliceID: sliceID}
	result, err := h.svc.View(c, sub)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func (h *HTTP) viewByName(c echo.Context) error {
	name := c.Param("name")
	if name == "" {
		return ErrMissingSubscriptionName
	}
	sub := sandpiper.Subscription{Name: name}
	result, err := h.svc.View(c, sub)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

// Subscription update request
type updateReq struct {
	SubID       uuid.UUID `json:"sub_id" validate:"required"`
	SliceID     uuid.UUID `json:"slice_id" validate:"required"`
	CompanyID   uuid.UUID `json:"company_id" validate:"required"`
	Name        string    `json:"name,omitempty" validate:"omitempty,min=3"`
	Description string    `json:"description,omitempty" validate:"omitempty"`
	Active      bool      `json:"active,omitempty" validate:"omitempty"`
}

func (h *HTTP) update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ErrInvalidSubscriptionUUID
	}
	req := new(updateReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	result, err := h.svc.Update(c, &subscription.Update{
		SubID:       id,
		SliceID:     req.SliceID,
		CompanyID:   req.CompanyID,
		Name:        req.Name,
		Description: req.Description,
		Active:      req.Active,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func (h *HTTP) delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ErrInvalidSubscriptionUUID
	}
	if err := h.svc.Delete(c, id); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
