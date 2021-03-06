// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package transport

// tag service routing functions

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/api/tag"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/params"
)

// HTTP represents user http service
type HTTP struct {
	svc tag.Service
}

// NewHTTP creates new tag http service
func NewHTTP(svc tag.Service, er *echo.Group) {
	h := HTTP{svc}

	er.POST("/tags", h.create)
	er.POST("/tags/:tagid/slices/:sliceid", h.assign)
	er.GET("/tags", h.list)
	er.GET("/tags/:id", h.view)
	er.PATCH("/tags/:id", h.update) // only update supplied fields
	er.DELETE("/tags/:id", h.delete)
	er.DELETE("/tags/:id/slices/:id", h.remove)
}

// Custom errors
var (
	// ErrNonNumericTagID indicates a malformed url parameter
	ErrNonNumericTagID  = echo.NewHTTPError(http.StatusBadRequest, "non-numeric tag id")
	ErrInvalidSliceUUID = echo.NewHTTPError(http.StatusBadRequest, "invalid slice uuid")
)

// Tag create request
type createReq struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// create populates createReq from supplied json body
func (h *HTTP) create(c echo.Context) error {
	r := new(createReq)
	if err := c.Bind(r); err != nil {
		return err
	}

	t := &sandpiper.Tag{
		Name:        r.Name,
		Description: r.Description,
	}

	// make sure valid tag name (stripping special chars)
	if err := t.CleanName(); err != nil {
		return err
	}

	result, err := h.svc.Create(c, *t)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, result)
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
	return c.JSON(http.StatusOK, sandpiper.TagsPaginated{Tags: result, Paging: p.Paging})
}

func (h *HTTP) view(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ErrNonNumericTagID
	}

	result, err := h.svc.View(c, id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

// Tag update request
type updateReq struct {
	ID          int    `json:"id" validate:"required"`
	Name        string `json:"name,omitempty" validate:"omitempty,min=3"`
	Description string `json:"description,omitempty" validate:"omitempty"`
}

func (h *HTTP) update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ErrNonNumericTagID
	}
	req := new(updateReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	result, err := h.svc.Update(c, &tag.Update{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func (h *HTTP) delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ErrNonNumericTagID
	}
	if err := h.svc.Delete(c, id); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

// assign adds a tag to a slice
func (h *HTTP) assign(c echo.Context) error {
	tagID, err := strconv.Atoi(c.Param("tagid"))
	if err != nil {
		return ErrNonNumericTagID
	}
	sliceID, err := uuid.Parse(c.Param("sliceid"))
	if err != nil {
		return ErrInvalidSliceUUID
	}
	if err := h.svc.Assign(c, tagID, sliceID); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

// remove deletes a tag from a slice
func (h *HTTP) remove(c echo.Context) error {
	tagID, err := strconv.Atoi(c.Param("tagid"))
	if err != nil {
		return ErrNonNumericTagID
	}
	sliceID, err := uuid.Parse(c.Param("sliceid"))
	if err != nil {
		return ErrInvalidSliceUUID
	}
	if err := h.svc.Remove(c, tagID, sliceID); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
