package minio

import (
	"aprilpollo/internal/adapters/config"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOClient struct {
	client *minio.Client
	cfg    *config.Minios3
}

func NewMinIOClient(cfg *config.Minios3) (*MinIOClient, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio: failed to create client: %w", err)
	}

	// ensure bucket exists
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("minio: bucket check failed: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("minio: failed to create bucket %q: %w", cfg.Bucket, err)
		}
	}

	return &MinIOClient{client: client, cfg: cfg}, nil
}

func (m *MinIOClient) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	_, err := m.client.PutObject(ctx, m.cfg.Bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("minio: upload failed: %w", err)
	}

	url := fmt.Sprintf("%s/%s/%s", m.cfg.EndpointPublic, m.cfg.Bucket, objectName)
	return url, nil
}

func (m *MinIOClient) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	u, err := m.client.PresignedGetObject(ctx, m.cfg.Bucket, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("minio: presign failed: %w", err)
	}
	return u.String(), nil
}

func (m *MinIOClient) DeleteFile(ctx context.Context, objectName string) error {
	err := m.client.RemoveObject(ctx, m.cfg.Bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("minio: delete failed: %w", err)
	}
	return nil
}
