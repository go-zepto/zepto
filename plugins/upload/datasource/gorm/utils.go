package gorm

import "github.com/mitchellh/mapstructure"

func decodeMapToStruct(input interface{}, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		Squash:           true,
		TagName:          "json",
		WeaklyTypedInput: true,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		panic(err)
	}
	return decoder.Decode(input)
}
