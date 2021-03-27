package linkeradmin

type FieldOptions = map[string]interface{}

type Field struct {
	Name    string       `json:"name"`
	Type    string       `json:"type"`
	Options FieldOptions `json:"options"`
}

type ResourceFieldEndpoint struct {
	Fields []*Field `json:"fields"`
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

func (e *ResourceFieldEndpoint) AddField(f Field) *ResourceFieldEndpoint {
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

func (e *ResourceFieldEndpoint) ReplaceField(name string, field Field) *ResourceFieldEndpoint {
	idx := e.findFieldIndexByName(name)
	if idx >= 0 {
		e.Fields[idx] = &field
	}
	return e
}

func NewResourceFieldEndpoint() *ResourceFieldEndpoint {
	return &ResourceFieldEndpoint{
		Fields: make([]*Field, 0),
	}
}
