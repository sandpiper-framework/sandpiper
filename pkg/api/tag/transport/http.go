// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package transport

// tag service routing functions

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"autocare.org/sandpiper/pkg/api/tag"
	"autocare.org/sandpiper/pkg/shared/model"
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
	er.PATCH("/tags/:id", h.update)
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
	return c.JSON(http.StatusOK, result)
}

type listResponse struct {
	Tags []sandpiper.Tag `json:"tags"`
	Page int             `json:"page"`
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
