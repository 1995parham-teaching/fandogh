package config

import (
	"github.com/1995parham/fandogh/internal/db"
	"github.com/1995parham/fandogh/internal/fs"
	"github.com/1995parham/fandogh/internal/logger"
	"github.com/1995parham/fandogh/internal/metric"
	telemetry "github.com/1995parham/fandogh/internal/telemetry/config"
)

// Default return default configuration.
func Default() Config {
	return Config{
		Logger: logger.Config{
			Level: "debug",
			Syslog: logger.Syslog{
				Enabled: false,
				Network: "",
				Address: "",
				Tag:     "",
			},
		},
		Database: db.Config{
			Name: "fandogh",
			URL:  "mongodb://127.0.0.1:27017",
		},
		FileStorage: fs.Config{
			Endpoint:  "127.0.0.1:9000",
			AccessKey: "access",
			SecretKey: "topsecret",
			UseSSL:    false,
		},
		Monitoring: metric.Config{
			Address: ":8080",
			Enabled: true,
		},
		Telemetry: telemetry.Config{
			Trace: telemetry.Trace{
				Enabled: false,
				URL:     "",
			},
		},
	}
}
