package broker

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"testing"
)

func TestEncodeJSONMessage(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	m := &Person{
		Name: "Ronaldo",
		Age:  43,
	}
	msg, err := encodeMessage(m)
	if err != nil {
		t.Fatal("encodeMessage should not throw error: ", msg)
	}
	var mm Person
	err = json.Unmarshal(msg.Body, &mm)
	if err != nil {
		t.Fatal("could not validate encoded json: ", err)
	}
	if mm.Name != "Ronaldo" || mm.Age != 43 {
		t.Fatal("could not validated encoded json. Received:", mm)
	}
}

func TestEncodeProtoMessage(t *testing.T) {
	p := wrapperspb.StringValue{Value: "Hello World"}
	msg, err := encodeMessage(&p)
	if err != nil {
		t.Fatal("encodeMessage should not throw error: ", msg)
	}
	var mm wrapperspb.StringValue
	err = proto.Unmarshal(msg.Body, &mm)
	if err != nil {
		t.Fatal("could not validate encoded proto message: ", err)
	}
	if mm.Value != "Hello World" {
		t.Fatal("could not validated encoded proto message. Received:", mm.Value)
	}
}
