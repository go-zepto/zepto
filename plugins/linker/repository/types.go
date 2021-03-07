package repository

import (
	"errors"
	"reflect"

	"github.com/go-zepto/zepto/plugins/linker/utils"
)

type SingleResult map[string]interface{}

func (s *SingleResult) Decode(dest interface{}) error {
	return utils.DecodeMapToStruct(s, dest)
}

type ManyResults []*SingleResult

func (m *ManyResults) Decode(dest interface{}) error {
	return utils.DecodeMapToStruct(m, dest)
}

type ListResult struct {
	Data  ManyResults `json:"data"`
	Count int64       `json:"count"`
}

func (s *ListResult) Decode(dest interface{}) error {
	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		return errors.New("decode dest should be a pointer")
	}
	if destType.Elem().Kind() == reflect.Slice {
		return utils.DecodeMapToStruct(s.Data, dest)
	}
	return utils.DecodeMapToStruct(s, dest)
}

type ManyAffectedResult struct {
	TotalAffected int64 `json:"total_affected"`
}

func (mar *ManyAffectedResult) Decode(dest interface{}) error {
	return utils.DecodeMapToStruct(mar, dest)
}
