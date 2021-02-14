package utils

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
)

// DecodeMapToStruct - Decode a map to struct.
func DecodeMapToStruct(input interface{}, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   output,
		TagName:  "json",
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		panic(err)
	}
	err = decoder.Decode(input)
	if err != nil {
		panic(err)
	}
	return decoder.Decode(input)
}

func DecodeStructToMap(input interface{}) map[string]interface{} {
	dest := map[string]interface{}{}
	s := reflect.Indirect(reflect.ValueOf(input))
	for i := 0; i < s.NumField(); i++ {
		valueField := s.Field(i)
		typeField := s.Type().Field(i)
		fieldName := ToSnakeCase(typeField.Name)
		jsonTag := typeField.Tag.Get("json")
		if jsonTag != "" {
			fieldName, _ = ParseJsonTag(jsonTag)
		}
		dest[fieldName] = valueField.Interface()
	}
	return dest
}
