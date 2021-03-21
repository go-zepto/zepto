package storage

import (
	"context"
	"io"
)

type UploadFileOptions struct {
	AccessType  AccessType
	Key         string
	Body        io.Reader
	ContentType string
}

type UploadFileResult struct {
	Key        string     `json:"key"`
	Url        string     `json:"url"`
	AccessType AccessType `json:"access_type"`
}

type DeleteFileOptions struct {
	Key string
}

type GenerateSignedURLOptions struct {
	Key string
}

type Storage interface {
	UploadFile(ctx context.Context, opts UploadFileOptions) (*UploadFileResult, error)
	DeleteFile(ctx context.Context, opts DeleteFileOptions) error
	GenerateSignedURL(ctx context.Context, opts GenerateSignedURLOptions) (string, error)
}
