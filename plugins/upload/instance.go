package upload

import (
	"context"

	"github.com/go-zepto/zepto/plugins/upload/datasource"
	"github.com/go-zepto/zepto/plugins/upload/storage"
	"github.com/go-zepto/zepto/web"
)

type UploadResponse struct {
	Key        string      `json:"key"`
	Url        string      `json:"url"`
	AccessType string      `json:"access_type"`
	File       interface{} `json:"file"`
}

type UploadInstance interface {
	UploadFile(ctx context.Context, opts storage.UploadFileOptions) (*UploadResponse, error)
	DeleteFile(ctx context.Context, opts storage.DeleteFileOptions) error
	GenerateSignedURL(ctx context.Context, opts storage.GenerateSignedURLOptions) (string, error)
}

type defaultUploadInstance struct {
	s  storage.Storage
	ds datasource.UploadDatasource
}

func (d *defaultUploadInstance) UploadFile(ctx context.Context, opts storage.UploadFileOptions) (*UploadResponse, error) {
	storageRes, err := d.s.UploadFile(ctx, opts)
	if err != nil {
		return nil, err
	}
	res := &UploadResponse{
		Key:        storageRes.Key,
		Url:        storageRes.Url,
		AccessType: storageRes.AccessType.String(),
	}
	// If the datasource is not configured, return the result (file is nil)
	if d.ds == nil {
		return res, nil
	}
	// Saving the file info using the Datasource
	file, err := d.ds.Create(ctx, &datasource.FileData{
		Key:        res.Key,
		Url:        res.Url,
		AccessType: res.AccessType,
	})
	if err != nil {
		return nil, err
	}
	res.File = file
	return res, nil
}

func (d *defaultUploadInstance) DeleteFile(ctx context.Context, opts storage.DeleteFileOptions) error {
	err := d.s.DeleteFile(ctx, opts)
	// If the datasource is not configured, return the result (file is nil)
	if d.ds == nil {
		return err
	}
	return d.ds.Delete(ctx, opts.Key)
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
