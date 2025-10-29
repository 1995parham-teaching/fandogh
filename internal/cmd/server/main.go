package server

import (
	"github.com/1995parham-teaching/fandogh/internal/config"
	"github.com/1995parham-teaching/fandogh/internal/db"
	"github.com/1995parham-teaching/fandogh/internal/fs"
	"github.com/1995parham-teaching/fandogh/internal/http/jwt"
	"github.com/1995parham-teaching/fandogh/internal/http/opa"
	"github.com/1995parham-teaching/fandogh/internal/http/server"
	"github.com/1995parham-teaching/fandogh/internal/logger"
	"github.com/1995parham-teaching/fandogh/internal/metric"
	"github.com/1995parham-teaching/fandogh/internal/store/home"
	"github.com/1995parham-teaching/fandogh/internal/store/user"
	"github.com/1995parham-teaching/fandogh/internal/telemetry/trace"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main(
	logger *zap.Logger,
	_ *echo.Echo,
	_ metric.Server,
) {
	logger.Info("welcome to fandogh server")
}

// Register server command.
func Register(root *cobra.Command) {
	root.AddCommand(
		// nolint: exhaustruct
		&cobra.Command{
			Use:   "server",
			Short: "Run server to serve the requests",
			Run: func(_ *cobra.Command, _ []string) {
				fx.New(
					fx.Provide(config.Provide),
					fx.Provide(logger.Provide),
					fx.Provide(trace.Provide),
					fx.Provide(db.Provide),
					fx.Provide(fs.Provide),
					fx.Provide(metric.Provide),
					fx.Provide(
						fx.Annotate(user.Provide, fx.As(new(user.User))),
					),
					fx.Provide(
						fx.Annotate(home.Provide, fx.As(new(home.Home))),
					),
					fx.Provide(jwt.Provide),
					fx.Provide(opa.Provide),
					fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
						return &fxevent.ZapLogger{Logger: logger}
					}),
					fx.Provide(server.Provide),
					fx.Invoke(main),
				).Run()
			},
		},
	)
}
