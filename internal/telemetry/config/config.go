package config

type Config struct {
	Trace `koanf:"trace"`
}

type Trace struct {
	Enabled bool   `koanf:"enabled"`
	Agent   string `koanf:"agent"`
}
