package server

import (
	"context"
	"net/http"

	"github.com/1995parham-teaching/fandogh/internal/http/handler"
	"github.com/1995parham-teaching/fandogh/internal/http/jwt"
	"github.com/1995parham-teaching/fandogh/internal/store/home"
	"github.com/1995parham-teaching/fandogh/internal/store/user"
	"github.com/labstack/echo/v5"
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
) *echo.Echo {
	app := echo.New()

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

	api := app.Group("/api", jwtHandler.Middleware())

	handler.Home{
		Store:  homeStore,
		Tracer: tracer,
		Logger: logger.Named("handler").Named("home"),
	}.Register(api)

	// nolint: exhaustruct
	server := &http.Server{
		Addr:    ":1378",
		Handler: app,
	}

	lc.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
				go func() {
					if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						logger.Fatal("echo initiation failed", zap.Error(err))
					}
				}()

				return nil
			},
			OnStop: server.Shutdown,
		},
	)

	return app
}
