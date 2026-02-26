package fs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Provide creates a connection to an S3-compatible storage (e.g., RustFS, MinIO).
func Provide(cfg Config) *s3.Client {
	scheme := "http"
	if cfg.UseSSL {
		scheme = "https"
	}

	// nolint: exhaustruct
	client := s3.New(s3.Options{
		BaseEndpoint: aws.String(fmt.Sprintf("%s://%s", scheme, cfg.Endpoint)),
		Region:       cfg.Region,
		Credentials:  credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		UsePathStyle: true,
	})

	return client
}

// Bucket ensures the specified bucket exists in the S3-compatible storage.
func Bucket(ctx context.Context, client *s3.Client, bucket string) error {
	// nolint: exhaustruct
	_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err == nil {
		return nil
	}

	// nolint: exhaustruct
	_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return fmt.Errorf("cannot create bucket [%s]: %w", bucket, err)
	}

	return nil
}
