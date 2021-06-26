package migrate

import (
	"github.com/1995parham/fandogh/internal/config"
	"github.com/1995parham/fandogh/internal/db"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const enable = 1

func main(cfg config.Config, logger *zap.Logger) {
	_, err := db.New(cfg.Database)
	if err != nil {
		logger.Fatal("database initiation failed", zap.Error(err))
	}
}

// Register migrate command.
func Register(root *cobra.Command, cfg config.Config, logger *zap.Logger) {
	root.AddCommand(
		// nolint: exhaustivestruct
		&cobra.Command{
			Use:   "migrate",
			Short: "Setup database indices",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg, logger)
			},
		},
	)
}
