package cmd

import (
	"os"

	"github.com/1995parham-teaching/fandogh/internal/cmd/migrate"
	"github.com/1995parham-teaching/fandogh/internal/cmd/server"
	"github.com/spf13/cobra"
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
	migrate.Register(root)

	if err := root.Execute(); err != nil {
		os.Exit(ExitFailure)
	}
}
