// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package transport

// auth routing functions

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/primary/auth"
)

// HTTP represents auth http service
type HTTP struct {
	svc auth.Service
}

// NewHTTP creates new auth http service
func NewHTTP(svc auth.Service, e *echo.Echo, mw echo.MiddlewareFunc) {
	h := HTTP{svc}
	e.POST("/login", h.login)
	e.GET("/refresh/:token", h.refresh)
	e.GET("/me", h.me, mw)
}

type credentials struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h *HTTP) login(c echo.Context) error {
	creds := new(credentials)
	if err := c.Bind(creds); err != nil {
		return err
	}
	r, err := h.svc.Authenticate(c, creds.Username, creds.Password)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, r)
}

func (h *HTTP) refresh(c echo.Context) error {
	r, err := h.svc.Refresh(c, c.Param("token"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, r)
}

func (h *HTTP) me(c echo.Context) error {
	user, err := h.svc.Me(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)
}
