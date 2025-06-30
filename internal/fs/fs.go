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
	// nolint: exhaustruct
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio connection failed: %w", err)
	}

	return client, nil
}

func Bucket(ctx context.Context, client *minio.Client, bucket string) error {
	found, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("cannot check minio bucket [%s] existence: %w", bucket, err)
	}

	if !found {
		// nolint: exhaustruct
		err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("cannot make minio bucket [%s]: %w", bucket, err)
		}
	}

	return nil
}
