package cmd

import (
	"os"

	"github.com/1995parham-teaching/fandogh/internal/cmd/migrate"
	"github.com/1995parham-teaching/fandogh/internal/cmd/server"
	"github.com/1995parham-teaching/fandogh/internal/config"
	"github.com/1995parham-teaching/fandogh/internal/logger"
	"github.com/1995parham-teaching/fandogh/internal/telemetry/trace"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// ExitFailure status code.
const ExitFailure = 1

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cfg := config.New()

	logger := logger.New(cfg.Logger)

	tracer := trace.New(cfg.Telemetry.Trace)

	// nolint: exhaustruct
	root := &cobra.Command{
		Use:   "fandogh",
		Short: "a server with login, user and etc.",
	}

	server.Register(root, cfg, logger, tracer)
	migrate.Register(root, cfg, logger)

	if err := root.Execute(); err != nil {
		logger.Error("failed to execute root command", zap.Error(err))
		os.Exit(ExitFailure)
	}
}
