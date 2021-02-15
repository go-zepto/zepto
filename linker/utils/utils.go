package utils

import (
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
)

func ToTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			return time.Parse(time.RFC3339, data.(string))
		case reflect.Float64:
			return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
		case reflect.Int64:
			return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
		default:
			return data, nil
		}
		// Convert it by parsing
	}
}

// DecodeMapToStruct - Decode a map to struct.
func DecodeMapToStruct(input interface{}, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   output,
		Squash:   true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			ToTimeHookFunc()),
		TagName: "json",
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		panic(err)
	}
	return decoder.Decode(input)
}

func DecodeStructToMap(input interface{}) map[string]interface{} {
	out := map[string]interface{}{}
	DecodeMapToStruct(input, &out)
	return out
}
