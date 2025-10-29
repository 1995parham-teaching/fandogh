package metric

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Config of metric servers.
type Config struct {
	Address string `koanf:"address"`
	Enabled bool   `koanf:"enabled"`
}

// Server contains information about metrics server.
type Server struct {
	srv     *http.Server
	enabled bool
}

// NewServer creates a new monitoring server.
func NewServer(cfg Config) Server {
	if !cfg.Enabled {
		return Server{
			srv:     nil,
			enabled: false,
		}
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	return Server{
		// nolint exhaustruct
		srv: &http.Server{
			Addr:              cfg.Address,
			Handler:           mux,
			ReadHeaderTimeout: time.Second,
		},
		enabled: true,
	}
}

// Provide creates a new monitoring server with lifecycle management.
func Provide(lc fx.Lifecycle, cfg Config, logger *zap.Logger) Server {
	server := NewServer(cfg)

	if !server.enabled {
		return server
	}

	lc.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
				go func() {
					logger.Info("starting metrics server", zap.String("address", server.srv.Addr))

					if err := server.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
						logger.Error("metric server initiation failed", zap.Error(err))
					}
				}()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				logger.Info("shutting down metrics server")

				return server.srv.Shutdown(ctx)
			},
		},
	)

	return server
}
