package output

import (
	"context"
	"io"
	"time"
)

type FileStorage interface {
	UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error)
	GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error)
	DeleteFile(ctx context.Context, objectName string) error
	DeleteFileByURL(ctx context.Context, fileURL string) error
	GetPresignedURLFromURL(ctx context.Context, fileURL string, expiry time.Duration) (string, error)
}
