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

// New creates a new mongodb connection and tests it.
func New(cfg Config) (*mongo.Database, error) {
	// connect to the mongodb
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

	// ping the mongodb
	{
		ctx, done := context.WithTimeout(context.Background(), connectionTimeout)
		defer done()

		err := client.Ping(ctx, readpref.Primary())
		if err != nil {
			return nil, fmt.Errorf("db ping error: %w", err)
		}
	}

	return client.Database(cfg.Name), nil
}

// Provide creates a new mongodb connection with lifecycle management.
func Provide(lc fx.Lifecycle, cfg Config, logger *zap.Logger) (*mongo.Database, error) {
	db, err := New(cfg)
	if err != nil {
		return nil, err
	}

	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				logger.Info("disconnecting from mongodb")
				return db.Client().Disconnect(ctx)
			},
		},
	)

	return db, nil
}
