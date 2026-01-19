package handler

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"

	"github.com/1995parham-teaching/fandogh/internal/http/common"
	intjwt "github.com/1995parham-teaching/fandogh/internal/http/jwt"
	"github.com/1995parham-teaching/fandogh/internal/http/request"
	"github.com/1995parham-teaching/fandogh/internal/model"
	"github.com/1995parham-teaching/fandogh/internal/store/home"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Home struct {
	Store  home.Home
	Tracer trace.Tracer
	Logger *zap.Logger
}

// New creates a home based on user request.
// Accepts JSON body with optional base64-encoded photos.
// nolint: wrapcheck, funlen, cyclop
func (h Home) New(c echo.Context) error {
	ctx, span := h.Tracer.Start(c.Request().Context(), "handler.home.create")
	defer span.End()

	var rq request.NewHome

	if err := c.Bind(&rq); err != nil {
		span.RecordError(err)

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := rq.Validate(); err != nil {
		span.RecordError(err)

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	token, ok := c.Get(common.UserContextKey).(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user claims not found")
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		span.RecordError(err)

		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	var bed model.Bed

	switch rq.Bed {
	case "single":
		bed = model.Single
	case "double":
		bed = model.Double
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "invalid bed type")
	}

	// Decode base64 photos
	photos := make([]model.Photo, 0, len(rq.Photos))

	for _, p := range rq.Photos {
		if p.Name == "" || p.Content == "" {
			continue
		}

		data, err := base64.StdEncoding.DecodeString(p.Content)
		if err != nil {
			span.RecordError(err)

			return echo.NewHTTPError(http.StatusBadRequest, "invalid base64 encoding for photo: "+p.Name)
		}

		photos = append(photos, model.Photo{
			Name:        p.Name,
			ContentType: http.DetectContentType(data),
			Content:     data,
		})
	}

	m := model.Home{
		ID:              "",
		Owner:           sub,
		Title:           rq.Title,
		Location:        rq.Location,
		Description:     rq.Description,
		Peoples:         rq.Peoples,
		Room:            rq.Room,
		Bed:             bed,
		Rooms:           rq.Rooms,
		Bathrooms:       rq.Bathrooms,
		Smoking:         rq.Smoking,
		Guest:           rq.Guest,
		Pet:             rq.Pet,
		BillsIncluded:   rq.BillsIncluded,
		Contract:        rq.Contract,
		SecurityDeposit: rq.SecurityDeposit,
		Photos:          nil,
		Price:           rq.Price,
	}

	if err := h.Store.Set(ctx, &m, photos); err != nil {
		span.RecordError(err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, m)
}

// Get retrieves a home by its ID.
// nolint: wrapcheck
func (h Home) Get(c echo.Context) error {
	ctx, span := h.Tracer.Start(c.Request().Context(), "handler.home.get")
	defer span.End()

	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "home id is required")
	}

	m, err := h.Store.Get(ctx, id)
	if err != nil {
		span.RecordError(err)

		if errors.Is(err, home.ErrIDNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, m)
}

// List retrieves homes with pagination.
// nolint: wrapcheck
func (h Home) List(c echo.Context) error {
	ctx, span := h.Tracer.Start(c.Request().Context(), "handler.home.list")
	defer span.End()

	skip := int64(0)
	limit := int64(10)

	if s := c.QueryParam("skip"); s != "" {
		if v, err := strconv.ParseInt(s, 10, 64); err == nil && v >= 0 {
			skip = v
		}
	}

	if l := c.QueryParam("limit"); l != "" {
		if v, err := strconv.ParseInt(l, 10, 64); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}

	result, err := h.Store.List(ctx, skip, limit)
	if err != nil {
		span.RecordError(err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

// Update modifies an existing home. Only the owner or an admin can update.
// nolint: wrapcheck, cyclop
func (h Home) Update(c echo.Context) error {
	ctx, span := h.Tracer.Start(c.Request().Context(), "handler.home.update")
	defer span.End()

	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "home id is required")
	}

	// Get the existing home to check ownership
	existingHome, err := h.Store.Get(ctx, id)
	if err != nil {
		span.RecordError(err)

		if errors.Is(err, home.ErrIDNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Get JWT claims
	token, ok := c.Get(common.UserContextKey).(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user claims not found")
	}

	claims, ok := token.Claims.(*intjwt.Claims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
	}

	sub, err := claims.GetSubject()
	if err != nil {
		span.RecordError(err)

		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	// Check authorization: must be owner or admin
	if existingHome.Owner != sub && !claims.Admin {
		return echo.NewHTTPError(http.StatusForbidden, "only the owner or an admin can update this home")
	}

	// Bind and validate request
	var rq request.UpdateHome

	if err := c.Bind(&rq); err != nil {
		span.RecordError(err)

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := rq.Validate(); err != nil {
		span.RecordError(err)

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var bed model.Bed

	switch rq.Bed {
	case "single":
		bed = model.Single
	case "double":
		bed = model.Double
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "invalid bed type")
	}

	updatedHome := model.Home{
		ID:              id,
		Owner:           existingHome.Owner,
		Title:           rq.Title,
		Location:        rq.Location,
		Description:     rq.Description,
		Peoples:         rq.Peoples,
		Room:            rq.Room,
		Bed:             bed,
		Rooms:           rq.Rooms,
		Bathrooms:       rq.Bathrooms,
		Smoking:         rq.Smoking,
		Guest:           rq.Guest,
		Pet:             rq.Pet,
		BillsIncluded:   rq.BillsIncluded,
		Contract:        rq.Contract,
		SecurityDeposit: rq.SecurityDeposit,
		Photos:          existingHome.Photos,
		Price:           rq.Price,
	}

	if err := h.Store.Update(ctx, id, updatedHome); err != nil {
		span.RecordError(err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, updatedHome)
}

func (h Home) Register(g *echo.Group) {
	g.POST("/homes", h.New)
	g.GET("/homes", h.List)
	g.GET("/homes/:id", h.Get)
	g.PUT("/homes/:id", h.Update)
}
