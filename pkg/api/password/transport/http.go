package transport

import (
	"net/http"
	"strconv"

	"autocare.org/sandpiper/pkg/api/password"

	"autocare.org/sandpiper/pkg/model"

	"github.com/labstack/echo/v4"
)

// HTTP represents password http transport service
type HTTP struct {
	svc password.Service
}

// NewHTTP creates new password http service
func NewHTTP(svc password.Service, er *echo.Group) {
	h := HTTP{svc}
	pr := er.Group("/password")
	pr.PATCH("/:id", h.change)
}

// Custom errors
var (
	ErrPasswordsNotMatching = echo.NewHTTPError(http.StatusBadRequest, "passwords do not match")
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
		return sandpiper.ErrBadRequest
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
