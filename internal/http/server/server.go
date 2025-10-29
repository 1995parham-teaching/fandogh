package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/1995parham-teaching/fandogh/internal/http/handler"
	"github.com/1995parham-teaching/fandogh/internal/http/jwt"
	"github.com/1995parham-teaching/fandogh/internal/http/opa"
	"github.com/1995parham-teaching/fandogh/internal/store/home"
	"github.com/1995parham-teaching/fandogh/internal/store/user"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Provide(
	lc fx.Lifecycle,
	userStore user.User,
	homeStore home.Home,
	logger *zap.Logger,
	tracer trace.Tracer,
	jwtHandler jwt.JWT,
	opaHandler opa.OPA,
) *echo.Echo {
	app := echo.New()
	app.Debug = true

	app.Use(otelecho.Middleware("fandogh"))

	handler.Healthz{
		Logger: logger.Named("handler").Named("healthz"),
		Tracer: tracer,
	}.Register(app.Group(""))

	handler.User{
		Store:  userStore,
		Tracer: tracer,
		Logger: logger.Named("handler").Named("user"),
		JWT:    jwtHandler,
	}.Register(app.Group(""))

	api := app.Group("/api", jwtHandler.Middleware(), opaHandler.Middleware())

	handler.Home{
		Store:  homeStore,
		Tracer: tracer,
		Logger: logger.Named("handler").Named("home"),
	}.Register(api)

	lc.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
				go func() {
					if err := app.Start(":1378"); !errors.Is(err, http.ErrServerClosed) {
						logger.Fatal("echo initiation failed", zap.Error(err))
					}
				}()

				return nil
			},
			OnStop: app.Shutdown,
		},
	)

	return app
}
