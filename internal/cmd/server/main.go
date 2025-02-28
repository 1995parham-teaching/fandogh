package server

import (
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/1995parham-teaching/fandogh/internal/config"
	"github.com/1995parham-teaching/fandogh/internal/db"
	"github.com/1995parham-teaching/fandogh/internal/fs"
	"github.com/1995parham-teaching/fandogh/internal/http/handler"
	"github.com/1995parham-teaching/fandogh/internal/http/jwt"
	"github.com/1995parham-teaching/fandogh/internal/metric"
	"github.com/1995parham-teaching/fandogh/internal/store/home"
	"github.com/1995parham-teaching/fandogh/internal/store/user"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func main(cfg config.Config, logger *zap.Logger, tracer trace.Tracer) {
	metric.NewServer(cfg.Monitoring).Start(logger.Named("metrics"))

	db, err := db.New(cfg.Database)
	if err != nil {
		logger.Fatal("database initiation failed", zap.Error(err))
	}

	mno, err := fs.New(cfg.FileStorage)
	if err != nil {
		logger.Fatal("file storage (minio) initiation failed", zap.Error(err))
	}

	app := echo.New()

	app.Use(otelecho.Middleware("fandogh"))

	handler.Healthz{
		Logger: logger.Named("handler").Named("healthz"),
		Tracer: tracer,
	}.Register(app.Group(""))

	jh := jwt.JWT{Config: cfg.JWT}

	handler.User{
		Store:  user.NewMongoUser(db, tracer),
		Tracer: tracer,
		Logger: logger.Named("handler").Named("user"),
		JWT:    jh,
	}.Register(app.Group(""))

	api := app.Group("/api", jh.Middleware())

	handler.Home{
		Store:  home.NewMongoHome(db, mno, tracer),
		Tracer: tracer,
		Logger: logger.Named("handler").Named("home"),
	}.Register(api)

	if err := app.Start(":1378"); !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal("echo initiation failed", zap.Error(err))
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

// Register server command.
func Register(root *cobra.Command, cfg config.Config, logger *zap.Logger, tracer trace.Tracer) {
	root.AddCommand(
		// nolint: exhaustruct
		&cobra.Command{
			Use:   "server",
			Short: "Run server to serve the requests",
			Run: func(_ *cobra.Command, _ []string) {
				main(cfg, logger, tracer)
			},
		},
	)
}
