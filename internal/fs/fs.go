package fs

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Provide creates a connection to an S3-compatible storage (e.g., SeaweedFS, MinIO).
func Provide(cfg Config) (*minio.Client, error) {
	// initialize S3-compatible client object.
	// nolint: exhaustruct
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("S3 storage connection failed: %w", err)
	}

	return client, nil
}

// Bucket ensures the specified bucket exists in the S3-compatible storage.
func Bucket(ctx context.Context, client *minio.Client, bucket string) error {
	found, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("cannot check bucket [%s] existence: %w", bucket, err)
	}

	if !found {
		// nolint: exhaustruct
		err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("cannot create bucket [%s]: %w", bucket, err)
		}
	}

	return nil
}
