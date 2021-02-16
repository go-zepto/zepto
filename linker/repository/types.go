package repository

import "github.com/go-zepto/zepto/linker/utils"

type SingleResult map[string]interface{}

func (s *SingleResult) Decode(dest interface{}) error {
	return utils.DecodeMapToStruct(s, dest)
}

type ListResult struct {
	Data  []map[string]interface{} `json:"data"`
	Count int64                    `json:"count"`
}

func (s *ListResult) Decode(dest interface{}) error {
	return utils.DecodeMapToStruct(s.Data, dest)
}

func (s *ListResult) DecodeAll(dest interface{}) error {
	return utils.DecodeMapToStruct(s, dest)
}

type ManyAffectedResult struct {
	TotalAffected int64
}
