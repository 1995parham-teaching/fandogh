package fs

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// New create a connection to minio to use as a file storage.
func New(cfg Config) (*minio.Client, error) {
	// initialize minio client object.
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio connection failed: %w", err)
	}

	return client, nil
}
