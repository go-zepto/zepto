package repository

import (
	"errors"
	"reflect"

	"github.com/go-zepto/zepto/linker/utils"
)

type SingleResult map[string]interface{}

func (s *SingleResult) Decode(dest interface{}) error {
	return utils.DecodeMapToStruct(s, dest)
}

type ListResult struct {
	Data  []map[string]interface{} `json:"data"`
	Count int64                    `json:"count"`
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
	TotalAffected int64
}

func (mar *ManyAffectedResult) Decode(dest interface{}) error {
	return utils.DecodeMapToStruct(mar, dest)
}
