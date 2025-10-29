package cmd

import (
	"os"

	"github.com/1995parham-teaching/fandogh/internal/cmd/migrate"
	"github.com/1995parham-teaching/fandogh/internal/cmd/server"
	"github.com/1995parham-teaching/fandogh/internal/config"
	"github.com/1995parham-teaching/fandogh/internal/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// ExitFailure status code.
const ExitFailure = 1

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// nolint: exhaustruct
	root := &cobra.Command{
		Use:   "fandogh",
		Short: "a server with login, user and etc.",
	}

	server.Register(root)

	// For migrate command, we still need config and logger
	cfg := config.Provide()
	logger := logger.Provide(cfg.Logger)
	migrate.Register(root, cfg, logger)

	err := root.Execute()
	if err != nil {
		if logger != nil {
			logger.Error("failed to execute root command", zap.Error(err))
		}
		os.Exit(ExitFailure)
	}
}
