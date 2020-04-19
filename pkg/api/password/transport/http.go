// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package transport

// password routing functions

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/password"
)

// HTTP represents password http transport service
type HTTP struct {
	svc password.Service
}

// NewHTTP creates new password http service
func NewHTTP(svc password.Service, er *echo.Group) {
	h := HTTP{svc}
	pr := er.Group("/password")
	pr.PATCH("/:id", h.change) // only changes provided fields
}

// Custom errors
var (
	ErrPasswordsNotMatching = echo.NewHTTPError(http.StatusBadRequest, "passwords do not match")
	ErrNonNumericID         = echo.NewHTTPError(http.StatusBadRequest, "numeric ID expected")
)

// Password change request
type changeReq struct {
	ID                 int    `json:"-"`
	OldPassword        string `json:"old_password" validate:"required,min=8"`
	NewPassword        string `json:"new_password" validate:"required,min=8"`
	NewPasswordConfirm string `json:"new_password_confirm" validate:"required"`
}

func (h *HTTP) change(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ErrNonNumericID
	}

	p := new(changeReq)
	if err := c.Bind(p); err != nil {
		return err
	}

	if p.NewPassword != p.NewPasswordConfirm {
		return ErrPasswordsNotMatching
	}

	if err := h.svc.Change(c, id, p.OldPassword, p.NewPassword); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
