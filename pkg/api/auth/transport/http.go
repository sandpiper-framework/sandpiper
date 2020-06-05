// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package transport

// auth routing functions

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"sandpiper/pkg/api/auth"
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
	e.GET("/server", h.server, mw)
}

func (h *HTTP) login(c echo.Context) error {
	creds, err := h.svc.ParseCredentials(c)
	if err != nil {
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

func (h *HTTP) server(c echo.Context) error {
	server := h.svc.Server(c)
	return c.JSON(http.StatusOK, server)
}
