package storage

import (
	"bytes"
	"encoding/json"
)

type AccessType int

const (
	Public AccessType = iota
	Private
)

var accessTypesStrings = [...]string{
	"public",
	"private",
}

var accessTypesToId = map[string]AccessType{
	"public":  Public,
	"private": Private,
}

func (a AccessType) String() string {
	return accessTypesStrings[a]
}

func (a AccessType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(a.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (a *AccessType) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}
	*a = accessTypesToId[j]
	return nil
}
