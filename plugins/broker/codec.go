package broker

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"reflect"
)

func encodeMessage(m interface{}) (*Message, error) {
	var msg Message
	pm, isProtoMessage := m.(proto.Message)
	// TODO: Create header from context
	header := make(map[string]string)
	if isProtoMessage {
		// Marshaling as proto message
		data, err := proto.Marshal(pm)
		if err != nil {
			return nil, err
		}
		msg = Message{
			Header: header,
			Body:   data,
		}
	} else {
		// Mashaling as JSON
		data, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}
		msg = Message{
			Header: header,
			Body:   data,
		}
	}
	return &msg, nil
}

func decodeMessage(m *Message, objType reflect.Type) (interface{}, error) {
	i := reflect.New(objType).Interface()
	pm, isProtoMessage := i.(proto.Message)
	if isProtoMessage {
		err := proto.Unmarshal(m.Body, pm)
		return pm, err
	}
	jm := i
	err := json.Unmarshal(m.Body, jm)
	return jm, err
}
