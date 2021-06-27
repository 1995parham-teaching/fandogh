package fs

type Config struct {
	Endpoint  string `koanf:"endpoint"`
	AccessKey string `koanf:"access_key"`
	SecretKey string `koanf:"secret_key"`
	UseSSL    bool   `koanf:"use_ssl"`
}
