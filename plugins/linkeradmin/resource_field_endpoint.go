package linkeradmin

import "github.com/go-zepto/zepto/plugins/linkeradmin/fields"

type ResourceFieldEndpoint struct {
	Fields []*fields.Field `json:"fields"`
}

func (e *ResourceFieldEndpoint) findFieldIndexByName(name string) int {
	for i := 0; i < len(e.Fields); i++ {
		if e.Fields[i].Name == name {
			return i
		}
	}
	return -1
}

func (e *ResourceFieldEndpoint) removeAtIndex(index int) *ResourceFieldEndpoint {
	e.Fields = append(e.Fields[:index], e.Fields[index+1:]...)
	return e
}

func (e *ResourceFieldEndpoint) AddField(f fields.Field) *ResourceFieldEndpoint {
	e.Fields = append(e.Fields, &f)
	return e
}

func (e *ResourceFieldEndpoint) RemoveField(name string) *ResourceFieldEndpoint {
	idx := e.findFieldIndexByName(name)
	if idx >= 0 {
		e.removeAtIndex(idx)
	}
	return e
}

func (e *ResourceFieldEndpoint) ReplaceField(name string, field fields.Field) *ResourceFieldEndpoint {
	idx := e.findFieldIndexByName(name)
	if idx >= 0 {
		e.Fields[idx] = &field
	}
	return e
}

func NewResourceFieldEndpoint() *ResourceFieldEndpoint {
	return &ResourceFieldEndpoint{
		Fields: make([]*fields.Field, 0),
	}
}
