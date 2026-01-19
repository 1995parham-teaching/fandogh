package migrate

import (
	"context"

	"github.com/1995parham-teaching/fandogh/internal/config"
	"github.com/1995parham-teaching/fandogh/internal/db"
	"github.com/1995parham-teaching/fandogh/internal/logger"
	"github.com/1995parham-teaching/fandogh/internal/store/user"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const enable = 1

func main(shutdowner fx.Shutdowner, logger *zap.Logger, db *mongo.Database) {
	idx, err := db.Collection(user.Collection).Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.M{"email": enable},
			Options: options.Index().SetUnique(true),
		})
	if err != nil {
		logger.Error("failed to create database index", zap.Error(err))
	}

	logger.Info("database index", zap.Any("index", idx))

	if err := shutdowner.Shutdown(); err != nil {
		logger.Error("failed to shutdown", zap.Error(err))
	}
}

// Register migrate command.
func Register(root *cobra.Command) {
	root.AddCommand(
		// nolint: exhaustruct
		&cobra.Command{
			Use:   "migrate",
			Short: "Setup database indices",
			Run: func(_ *cobra.Command, _ []string) {
				fx.New(
					fx.Provide(config.Provide),
					fx.Provide(logger.Provide),
					fx.Provide(db.Provide),
					fx.Options(fx.NopLogger),
					fx.Invoke(main),
				).Run()
			},
		},
	)
}
