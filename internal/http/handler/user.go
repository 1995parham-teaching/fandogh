package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/1995parham/fandogh/internal/http/jwt"
	"github.com/1995parham/fandogh/internal/http/request"
	"github.com/1995parham/fandogh/internal/http/response"
	"github.com/1995parham/fandogh/internal/model"
	"github.com/1995parham/fandogh/internal/store/user"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type User struct {
	Store  user.User
	Tracer trace.Tracer
	Logger *zap.Logger
	JWT    jwt.JWT
}

// nolint: wrapcheck
func (h User) Create(c echo.Context) error {
	ctx, span := h.Tracer.Start(c.Request().Context(), "handler.user.create")
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

// Login checks given credentials and generate jwt token
// nolint: wrapcheck
func (h User) Login(c echo.Context) error {
	ctx, span := h.Tracer.Start(c.Request().Context(), "handler.user.login")
	defer span.End()

	var rq request.Login

	if err := c.Bind(&rq); err != nil {
		span.RecordError(err)

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := rq.Validate(); err != nil {
		span.RecordError(err)

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	u, err := h.Store.Get(ctx, rq.Email)
	if err != nil {
		span.RecordError(err)

		if errors.Is(err, user.ErrEmailNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("email %s does not exist", rq.Email))
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if u.Password != rq.Password {
		return echo.NewHTTPError(http.StatusUnauthorized, "incorrect password")
	}

	var res response.Login

	res.User = u

	t, err := h.JWT.NewAccessToken(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res.AccessToken = t

	return c.JSON(http.StatusOK, res)
}

// Register registers the routes of User handler on given group.
func (h User) Register(g *echo.Group) {
	g.POST("/register", h.Create)
	g.POST("/login", h.Login)
}
