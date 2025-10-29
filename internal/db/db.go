package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/v2/mongo/otelmongo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const connectionTimeout = 10 * time.Second

// Provide creates a new mongodb connection with lifecycle management.
func Provide(lc fx.Lifecycle, cfg Config, logger *zap.Logger) (*mongo.Database, error) {
	opts := options.Client()
	opts.ApplyURI(cfg.URL)
	opts.SetMonitor(otelmongo.NewMonitor())
	opts.SetConnectTimeout(connectionTimeout)

	client, err := mongo.Connect(
		opts,
	)
	if err != nil {
		return nil, fmt.Errorf("db connection error: %w", err)
	}

	db := client.Database(cfg.Name)

	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				logger.Info("disconnecting from mongodb")

				return db.Client().Disconnect(ctx)
			},
			OnStart: func(ctx context.Context) error {
				err := client.Ping(ctx, readpref.Primary())
				if err != nil {
					return fmt.Errorf("db ping error: %w", err)
				}

				return nil
			},
		},
	)

	return db, nil
}
