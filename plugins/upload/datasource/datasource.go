package datasource

import "context"

type FileData struct {
	Key        string `json:"key"`
	Url        string `json:"url"`
	AccessType string `json:"access_type"`
}

type UploadDatasource interface {
	Create(ctx context.Context, data *FileData) (interface{}, error)
	Delete(ctx context.Context, key string) error
}
