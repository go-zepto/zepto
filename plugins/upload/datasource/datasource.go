package datasource

type FileData struct {
	Key        string `json:"key"`
	Url        string `json:"url"`
	AccessType string `json:"access_type"`
}

type UploadDatasource interface {
	Create(data *FileData) (interface{}, error)
	Delete(key string) error
}
