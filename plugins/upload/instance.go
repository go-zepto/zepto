package upload

import (
	"context"

	"github.com/go-zepto/zepto/plugins/upload/storage"
	"github.com/go-zepto/zepto/web"
)

type UploadInstance interface {
	storage.Storage
}

type defaultUploadInstance struct {
	s storage.Storage
}

func (d *defaultUploadInstance) UploadFile(ctx context.Context, opts storage.UploadFileOptions) (*storage.UploadFileResult, error) {
	return d.s.UploadFile(ctx, opts)
}

func (d *defaultUploadInstance) DeleteFile(ctx context.Context, opts storage.DeleteFileOptions) error {
	return d.s.DeleteFile(ctx, opts)
}

func (d *defaultUploadInstance) GenerateSignedURL(ctx context.Context, opts storage.GenerateSignedURLOptions) (string, error) {
	return d.s.GenerateSignedURL(ctx, opts)
}

func InstanceFromCtx(ctx web.Context) UploadInstance {
	i := ctx.PluginInstance("upload")
	uploadInstance, ok := i.(UploadInstance)
	if !ok {
		return nil
	}
	return uploadInstance
}
