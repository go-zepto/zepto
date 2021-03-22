package models

type File struct {
	Model
	Key        string `json:"key"`
	Url        string `json:"url"`
	AccessType string `json:"access_type"`
}
