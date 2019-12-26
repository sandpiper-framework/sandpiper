package transport

// slice routing functions

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/satori/go.uuid"

	"autocare.org/sandpiper/internal/model"
	"autocare.org/sandpiper/pkg/api/slice"
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
	sr.GET("", h.list)
	sr.GET("/:id", h.view)
	sr.PATCH("/:id", h.update)
	sr.DELETE("/:id", h.delete)
}

// Custom errors
var (
	ErrUnknownRoleID    = echo.NewHTTPError(http.StatusBadRequest, "unknown access role")
	ErrInvalidSliceUUID = echo.NewHTTPError(http.StatusBadRequest, "invalid slice uuid")
)

// Slice create request
type createReq struct {
	Name         string    `json:"name" validate:"required,min=3"`
	ContentHash  string    `json:"content_hash"`
	ContentCount uint      `json:"content_count"`
	LastUpdate   time.Time `json:"last_update"`

	CompanyID  int                  `json:"company_id" validate:"required"`
	LocationID int                  `json:"location_id" validate:"required"`
	RoleID     sandpiper.AccessRole `json:"role_id" validate:"required"`
}

func (h *HTTP) create(c echo.Context) error {
	r := new(createReq)

	if err := c.Bind(r); err != nil {
		return err
	}

	if r.RoleID < sandpiper.SuperAdminRole || r.RoleID > sandpiper.UserRole {
		return ErrUnknownRoleID
	}

	usr, err := h.svc.Create(c, sandpiper.Slice{
		Name:         r.Name,
		ContentHash:  r.ContentHash,
		ContentCount: r.ContentCount,
		LastUpdate:   r.LastUpdate,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, usr)
}

type listResponse struct {
	Slices []sandpiper.Slice `json:"slices"`
	Page   int               `json:"page"`
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
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return ErrInvalidSliceUUID
	}

	result, err := h.svc.View(c, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

// Slice update request
type updateReq struct {
	ID           uuid.UUID `json:"-"`
	Name         string    `json:"name,omitempty" validate:"omitempty,min=3"`
	ContentHash  string    `json:"content_hash,omitempty" validate:"omitempty,min=2"`
	ContentCount uint      `json:"content_count,omitempty"`
	LastUpdate   time.Time `json:"last_update,omitempty"`
}

func (h *HTTP) update(c echo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return ErrInvalidSliceUUID
	}

	req := new(updateReq)
	if err := c.Bind(req); err != nil {
		return err
	}

	usr, err := h.svc.Update(c, &slice.Update{
		ID:           id,
		Name:         req.Name,
		ContentHash:  req.ContentHash,
		ContentCount: req.ContentCount,
		LastUpdate:   req.LastUpdate,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, usr)
}

func (h *HTTP) delete(c echo.Context) error {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return ErrInvalidSliceUUID
	}

	if err := h.svc.Delete(c, id); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
