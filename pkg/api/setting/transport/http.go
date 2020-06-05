// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package transport

// setting routing

// database settings are stored in a one-row table (enforced at the database level)

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"sandpiper/pkg/api/setting"
	"sandpiper/pkg/shared/model"
)

// HTTP represents user http service
type HTTP struct {
	svc setting.Service
}

// NewHTTP creates new setting http service
func NewHTTP(svc setting.Service, er *echo.Group) {
	h := HTTP{svc}
	sr := er.Group("/settings")
	sr.POST("", h.create)
	sr.GET("", h.view)
	er.PUT("", h.update) // not a PATCH, body must include *all* fields
}

// Setting create request
type createReq struct {
	ID         bool      `json:"id"`
	ServerRole string    `json:"server_role" validate:"required"`
	ServerID   uuid.UUID `json:"server_id" validate:"required"`
}

func (h *HTTP) create(c echo.Context) error {
	r := new(createReq)
	if err := c.Bind(r); err != nil {
		return err
	}

	result, err := h.svc.Create(c, &sandpiper.Setting{
		ID:         true,
		ServerRole: r.ServerRole,
		ServerID:   r.ServerID,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, result)
}

func (h *HTTP) view(c echo.Context) error {
	result, err := h.svc.View(c)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

// Subscription update request
type updateReq struct {
	ID         bool      `json:"id"`
	ServerRole string    `json:"server_role" validate:"required"`
	ServerID   uuid.UUID `json:"server_id" validate:"required"`
}

func (h *HTTP) update(c echo.Context) error {
	req := new(updateReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	result, err := h.svc.Update(c, &setting.Update{
		ID:         true,
		ServerRole: req.ServerRole,
		ServerID:   req.ServerID,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}
