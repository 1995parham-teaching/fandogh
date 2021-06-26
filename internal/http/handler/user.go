package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/1995parham/fandogh/internal/http/request"
	"github.com/1995parham/fandogh/internal/model"
	"github.com/1995parham/fandogh/internal/store/user"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
)

type User struct {
	Store  user.User
	Tracer trace.Tracer
}

// nolint: wrapcheck
func (h User) Create(c echo.Context) error {
	ctx, span := h.Tracer.Start(c.Request().Context(), "handler.url.create")
	defer span.End()

	var rq request.Register

	if err := c.Bind(&rq); err != nil {
		span.RecordError(err)

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := rq.Validate(); err != nil {
		span.RecordError(err)

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	u := model.User{
		Email:    rq.Email,
		Password: rq.Password,
		Name:     rq.Name,
	}

	if err := h.Store.Set(ctx, u); err != nil {
		span.RecordError(err)

		if errors.Is(err, user.ErrEmailDuplicate) {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("email %s already exists", u.Email))
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, u)
}

// Register registers the routes of User handler on given group.
func (h User) Register(g *echo.Group) {
	g.POST("/register", h.Create)
}
