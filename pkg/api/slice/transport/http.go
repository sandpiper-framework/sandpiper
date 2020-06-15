// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

package transport

// slice routing functions

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/sandpiper-framework/sandpiper/pkg/api/slice"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/params"
)

// HTTP represents user http service
type HTTP struct {
	svc slice.Service
}

// NewHTTP creates new slice http service
func NewHTTP(svc slice.Service, er *echo.Group) {
	h := HTTP{svc}
	sr := er.Group("/slices")
	sr.POST("", h.create)
	sr.POST("/refresh/:id", h.refresh)
	sr.GET("", h.list) // filter by tags (/slices?tags=aaa,bbb or /slices?tags-all=aaa,bbb)
	sr.GET("/:id", h.view)
	sr.PATCH("/:id", h.update) // only update supplied fields
	sr.DELETE("/:id", h.delete)
	sr.GET("/name/:name", h.viewByName)
	sr.GET("/metadata/:id", h.metadata)
	sr.PUT("/lock/:id", h.lock)
	sr.PUT("/unlock/:id", h.unlock)
}

// Custom errors
var (
	// ErrInvalidSliceUUID indicates a malformed uuid
	ErrInvalidSliceUUID = echo.NewHTTPError(http.StatusBadRequest, "invalid slice uuid")
	ErrInvalidSliceType = echo.NewHTTPError(http.StatusBadRequest, "invalid slice-type")
)

// Slice create request
type createReq struct {
	ID        uuid.UUID         `json:"id"` // optional
	Name      string            `json:"name" validate:"required,min=3"`
	SliceType string            `json:"slice_type" validate:"required"`
	AllowSync bool              `json:"allow_sync"`
	Metadata  sandpiper.MetaMap `json:"metadata"`
}

func (r createReq) id() uuid.UUID {
	if r.ID == uuid.Nil {
		return uuid.New()
	}
	return r.ID
}

// create populates createReq from body json adding UUID if not provided
func (h *HTTP) create(c echo.Context) error {
	r := new(createReq)
	if err := c.Bind(r); err != nil {
		return err
	}

	rec := sandpiper.Slice{
		ID:           r.id(),
		Name:         r.Name,
		SliceType:    r.SliceType,
		AllowSync:    r.AllowSync,
		SyncStatus:   sandpiper.SyncStatusNone,
		ContentHash:  "",
		ContentCount: 0,
		Metadata:     r.Metadata,
	}

	if !rec.Validate() {
		return ErrInvalidSliceType
	}

	result, err := h.svc.Create(c, rec)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, result)
}

// refresh content information after slice grains are changed
func (h *HTTP) refresh(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ErrInvalidSliceUUID
	}

	if err := h.svc.Refresh(c, id); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *HTTP) view(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ErrInvalidSliceUUID
	}

	result, err := h.svc.View(c, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func (h *HTTP) viewByName(c echo.Context) error {
	result, err := h.svc.ViewByName(c, c.Param("name"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func (h *HTTP) list(c echo.Context) error {
	// allow slices filtered by tags (/slices?tags=aaa,bbb or /slices?tags-all=aaa,bbb)
	tags := params.NewTagQuery(c.QueryParams())

	p, err := params.Parse(c)
	if err != nil {
		return err
	}

	result, err := h.svc.List(c, p, tags)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, sandpiper.SlicesPaginated{Slices: result, Paging: p.Paging})
}

func (h *HTTP) metadata(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ErrInvalidSliceUUID
	}

	result, err := h.svc.Metadata(c, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

// Slice update request
type updateReq struct {
	ID           uuid.UUID `json:"-"`
	Name         string    `json:"name,omitempty" validate:"omitempty,min=3"`
	SLiceType    string    `json:"slice_type" validate:"required"`
	ContentHash  string    `json:"content_hash,omitempty" validate:"omitempty,min=2"`
	ContentCount int       `json:"content_count,omitempty"`
	ContentDate  time.Time `json:"content_date,omitempty"`
	AllowSync    bool      `json:"allow_sync"`
}

func (h *HTTP) update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ErrInvalidSliceUUID
	}

	req := new(updateReq)
	if err := c.Bind(req); err != nil {
		return err
	}

	result, err := h.svc.Update(c, &slice.Update{
		ID:           id,
		Name:         req.Name,
		SliceType:    req.SLiceType,
		ContentHash:  req.ContentHash,
		ContentCount: req.ContentCount,
		ContentDate:  req.ContentDate,
		AllowSync:    req.AllowSync,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func (h *HTTP) delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ErrInvalidSliceUUID
	}

	if err := h.svc.Delete(c, id); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *HTTP) lock(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ErrInvalidSliceUUID
	}

	if err := h.svc.Lock(c, id); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *HTTP) unlock(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ErrInvalidSliceUUID
	}

	if err := h.svc.Unlock(c, id); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
