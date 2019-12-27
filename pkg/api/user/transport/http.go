// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package transport

// user routing functions

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/satori/go.uuid"

	"autocare.org/sandpiper/internal/model"
	"autocare.org/sandpiper/pkg/api/user"
)

// HTTP represents user http service
type HTTP struct {
	svc user.Service
}

// NewHTTP creates new user http service
func NewHTTP(svc user.Service, er *echo.Group) {
	h := HTTP{svc}
	ur := er.Group("/users")
	ur.POST("", h.create)
	ur.GET("", h.list)
	ur.GET("/:id", h.view)
	ur.PATCH("/:id", h.update)
	ur.DELETE("/:id", h.delete)
}

// Custom errors
var (
	ErrPasswordsNotMatching = echo.NewHTTPError(http.StatusBadRequest, "passwords do not match")
	ErrUnknownRoleID        = echo.NewHTTPError(http.StatusBadRequest, "unknown access role")
	ErrNonNumericUserID     = echo.NewHTTPError(http.StatusBadRequest, "numeric user id expected")
)

// User create request
type createReq struct {
	FirstName       string `json:"first_name" validate:"required"`
	LastName        string `json:"last_name" validate:"required"`
	Username        string `json:"username" validate:"required,min=3,alphanum"`
	Password        string `json:"password" validate:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" validate:"required"`
	Email           string `json:"email" validate:"required,email"`

	CompanyID uuid.UUID            `json:"company_id" validate:"required"`
	RoleID    sandpiper.AccessRole `json:"role_id" validate:"required"`
}

func (h *HTTP) create(c echo.Context) error {
	r := new(createReq)

	if err := c.Bind(r); err != nil {
		return err
	}

	if r.Password != r.PasswordConfirm {
		return ErrPasswordsNotMatching
	}

	if r.RoleID < sandpiper.SuperAdminRole || r.RoleID > sandpiper.UserRole {
		return ErrUnknownRoleID
	}

	usr, err := h.svc.Create(c, sandpiper.User{
		Username:  r.Username,
		Password:  r.Password,
		Email:     r.Email,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		CompanyID: r.CompanyID,
		RoleID:    r.RoleID,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, usr)
}

type listResponse struct {
	Users []sandpiper.User `json:"users"`
	Page  int              `json:"page"`
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ErrNonNumericUserID
	}

	result, err := h.svc.View(c, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

// User update request
type updateReq struct {
	ID        int    `json:"-"`
	FirstName string `json:"first_name,omitempty" validate:"omitempty,min=2"`
	LastName  string `json:"last_name,omitempty" validate:"omitempty,min=2"`
	Mobile    string `json:"mobile,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Address   string `json:"address,omitempty"`
}

func (h *HTTP) update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ErrNonNumericUserID
	}

	req := new(updateReq)
	if err := c.Bind(req); err != nil {
		return err
	}

	usr, err := h.svc.Update(c, &user.Update{
		ID:        id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Mobile:    req.Mobile,
		Phone:     req.Phone,
		Address:   req.Address,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, usr)
}

func (h *HTTP) delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ErrNonNumericUserID
	}

	if err := h.svc.Delete(c, id); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
