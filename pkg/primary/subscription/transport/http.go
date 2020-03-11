// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package transport

// subscription service routing functions

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/primary/subscription"
	"autocare.org/sandpiper/pkg/shared/model"
)

// HTTP represents user http service
type HTTP struct {
	svc subscription.Service
}

// NewHTTP creates new subscription http service
func NewHTTP(svc subscription.Service, er *echo.Group) {
	h := HTTP{svc}

	er.GET("/companies/:id/subscriptions", h.listByCompany)
	er.GET("/companies/:id/subscriptions/:sliceid", h.viewByCompany)

	er.POST("/subscriptions", h.create)
	er.GET("/subscriptions", h.list)
	er.GET("/subscriptions/:id", h.view)
	er.PATCH("/subscriptions/:id", h.update)
	er.DELETE("/subscriptions/:id", h.delete)
}

// Custom errors
var (
	ErrInvalidCompanyUUID       = echo.NewHTTPError(http.StatusBadRequest, "malformed company uuid")
	ErrInvalidSliceUUID         = echo.NewHTTPError(http.StatusBadRequest, "malformed slice uuid")
	ErrNonNumericSubscriptionID = echo.NewHTTPError(http.StatusBadRequest, "non-numeric subscription id")
)

// Subscription create request
type createReq struct {
	SliceID     uuid.UUID `json:"slice_id" validate:"required"`
	CompanyID   uuid.UUID `json:"company_id" validate:"required"`
	Name        string    `json:"name" validate:"required,min=3"`
	Description string    `json:"description"`
	Active      bool      `json:"active"`
}

// create populates createReq from supplied json body
func (h *HTTP) create(c echo.Context) error {
	r := new(createReq)
	if err := c.Bind(r); err != nil {
		return err
	}
	result, err := h.svc.Create(c, sandpiper.Subscription{
		SliceID:     r.SliceID,
		CompanyID:   r.CompanyID,
		Name:        r.Name,
		Description: r.Description,
		Active:      r.Active,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

type listResponse struct {
	Subscriptions []sandpiper.Subscription `json:"subscriptions"`
	Page          int                      `json:"page"`
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
	return c.JSON(http.StatusOK, listResponse{result, p.Page})
}

func (h *HTTP) view(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ErrNonNumericSubscriptionID
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

// Subscription update request
type updateReq struct {
	SubID int `json:"sub_id" validate:"required"`
	//SliceID     uuid.UUID `json:"slice_id" validate:"required"`
	//CompanyID   uuid.UUID `json:"company_id" validate:"required"`
	Name        string `json:"name,omitempty" validate:"omitempty,min=3"`
	Description string `json:"description,omitempty" validate:"omitempty"`
	Active      bool   `json:"active,omitempty" validate:"omitempty"`
}

func (h *HTTP) update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ErrNonNumericSubscriptionID
	}
	req := new(updateReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	result, err := h.svc.Update(c, &subscription.Update{
		SubID: id,
		//SliceID:     req.SliceID,
		//CompanyID:   req.CompanyID,
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ErrNonNumericSubscriptionID
	}
	if err := h.svc.Delete(c, id); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
