package storage

import (
	"context"
	"io"
	"time"
)

type CreateOptions struct {
	AccessType  AccessType
	Key         string
	Body        io.Reader
	ContentType string
	Size        int64
}

type UpdateOptions struct {
	Key string
}

type FileInfo struct {
	AccessType  AccessType
	Key         string
	ContentType string
	Size        int64
	Url         string
}

type Storage interface {
	Create(ctx context.Context, opts CreateOptions) (*FileInfo, error)
	Destroy(ctx context.Context, key string) error
	GenerateSignedURL(ctx context.Context, key string, expirationTime time.Duration) (string, error)
}
