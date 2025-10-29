package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1995parham-teaching/fandogh/internal/http/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

type HealthzSuite struct {
	suite.Suite

	engine *echo.Echo
}

func (suite *HealthzSuite) SetupSuite() {
	var engine *echo.Echo

	app := fxtest.New(
		suite.T(),
		fx.Provide(func() *zap.Logger {
			return zap.NewNop()
		}),
		fx.Provide(func() trace.Tracer {
			return noop.NewTracerProvider().Tracer("")
		}),
		fx.Provide(func(logger *zap.Logger, tracer trace.Tracer) *echo.Echo {
			e := echo.New()
			handler.Healthz{
				Logger: logger,
				Tracer: tracer,
			}.Register(e.Group(""))
			return e
		}),
		fx.Populate(&engine),
	)
	defer app.RequireStart().RequireStop()

	suite.engine = engine
}

func (suite *HealthzSuite) TestHandler() {
	require := suite.Require()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	suite.engine.ServeHTTP(w, req)
	require.Equal(http.StatusNoContent, w.Code)
}

func TestHealthzSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HealthzSuite))
}
