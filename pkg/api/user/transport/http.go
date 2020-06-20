// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package transport

// user routing functions

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/api/user"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/params"
)

// HTTP represents user http service
type HTTP struct {
	svc user.Service
}

// NewHTTP creates new user http service
func NewHTTP(svc user.Service, er *echo.Group) {
	h := HTTP{svc}
	er.POST("/apikey", h.createAPIKey)
	ur := er.Group("/users")
	ur.POST("", h.create)
	ur.GET("", h.list)
	ur.GET("/:id", h.view)
	ur.PATCH("/:id", h.update) // only update supplied fields
	ur.DELETE("/:id", h.delete)
}

// Custom errors
var (
	ErrPasswordsNotMatching = echo.NewHTTPError(http.StatusBadRequest, "Passwords do not match")
	ErrUnknownRole          = echo.NewHTTPError(http.StatusBadRequest, "Unknown access role")
	ErrNonNumericUserID     = echo.NewHTTPError(http.StatusBadRequest, "Numeric user id expected")
)

// User create request (id is assigned by database)
type createReq struct {
	FirstName       string                `json:"first_name" validate:"required"`
	LastName        string                `json:"last_name" validate:"required"`
	Username        string                `json:"username" validate:"required,min=3,alphanum"`
	Password        string                `json:"password" validate:"required,min=8"`
	PasswordConfirm string                `json:"password_confirm" validate:"required"`
	Email           string                `json:"email" validate:"required,email"`
	CompanyID       uuid.UUID             `json:"company_id" validate:"required"`
	Role            sandpiper.AccessLevel `json:"role" validate:"required"`
	Active          bool                  `json:"active"`
}

func (h *HTTP) create(c echo.Context) error {
	r := new(createReq)

	if err := c.Bind(r); err != nil {
		return err
	}

	if r.Password != r.PasswordConfirm {
		return ErrPasswordsNotMatching
	}

	if !sandpiper.RoleIsValid(r.Role) {
		return ErrUnknownRole
	}

	usr, err := h.svc.Create(c, sandpiper.User{
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Username:  r.Username,
		Password:  r.Password,
		Email:     r.Email,
		CompanyID: r.CompanyID,
		Role:      r.Role,
		Active:    r.Active,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, usr)
}

func (h *HTTP) list(c echo.Context) error {
	p, err := params.Parse(c)
	if err != nil {
		return err
	}

	result, err := h.svc.List(c, p)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, sandpiper.UsersPaginated{Users: result, Paging: p.Paging})
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

// User update request (only things they can change here)
type updateReq struct {
	ID        int    `json:"-"`
	FirstName string `json:"first_name,omitempty" validate:"omitempty,min=2"`
	LastName  string `json:"last_name,omitempty" validate:"omitempty,min=2"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Active    bool   `json:"active,omitempty"`
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
		Email:     req.Email,
		Phone:     req.Phone,
		Active:    req.Active,
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

func (h *HTTP) createAPIKey(c echo.Context) error {
	result, err := h.svc.CreateAPIKey(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}
