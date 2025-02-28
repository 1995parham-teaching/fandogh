package migrate

import (
	"context"

	"github.com/1995parham-teaching/fandogh/internal/config"
	"github.com/1995parham-teaching/fandogh/internal/db"
	"github.com/1995parham-teaching/fandogh/internal/store/user"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const enable = 1

func main(cfg config.Config, logger *zap.Logger) {
	db, err := db.New(cfg.Database)
	if err != nil {
		logger.Fatal("database initiation failed", zap.Error(err))
	}

	idx, err := db.Collection(user.Collection).Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.M{"email": enable},
			Options: options.Index().SetUnique(true),
		})
	if err != nil {
		panic(err)
	}

	logger.Info("database index", zap.Any("index", idx))
}

// Register migrate command.
func Register(root *cobra.Command, cfg config.Config, logger *zap.Logger) {
	root.AddCommand(
		// nolint: exhaustruct
		&cobra.Command{
			Use:   "migrate",
			Short: "Setup database indices",
			Run: func(_ *cobra.Command, _ []string) {
				main(cfg, logger)
			},
		},
	)
}
