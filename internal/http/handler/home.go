package handler

import (
	"io"
	"net/http"

	"github.com/1995parham-teaching/fandogh/internal/http/common"
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

// nolint: wrapcheck, funlen, cyclop
// New create a home based on user request. unlike the other requests this request uses form.
// this request can contains many files as home images.
// please note that the maximum size of each home creation request is 10MB.
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

	mc, ok := c.Get(common.UserContextKey).(*jwt.Token)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "user claims not found")
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

	var photoFields []string
	if err := echo.FormFieldBinder(c).BindWithDelimiter("photos", &photoFields, ",").BindError(); err != nil {
		span.RecordError(err)

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	photos := make([]model.Photo, len(photoFields))

	for i, field := range photoFields {
		image, err := c.FormFile(field)
		if err != nil {
			span.RecordError(err)

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		fd, err := image.Open()
		if err != nil {
			span.RecordError(err)

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		defer fd.Close()

		data, err := io.ReadAll(fd)
		if err != nil {
			span.RecordError(err)

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		photos[i] = model.Photo{
			Name:        field,
			ContentType: http.DetectContentType(data),
			Content:     data,
		}
	}

	// Get the subject from the token
	sub, err := mc.Claims.GetSubject()
	if err != nil {
		span.RecordError(err)

		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
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

func (h Home) Register(g *echo.Group) {
	g.POST("/homes", h.New)
}
