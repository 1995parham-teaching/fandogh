package fs

import (
	"context"
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

	found, err := client.BucketExists(context.Background(), cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("cannot check minio bucket [%s] existence: %w", cfg.Bucket, err)
	}

	if !found {
		if err := client.MakeBucket(context.Background(), cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("cannot make minio bucket [%s]: %w", cfg.Bucket, err)
		}
	}

	return client, nil
}
