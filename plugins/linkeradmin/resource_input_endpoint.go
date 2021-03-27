package linkeradmin

// Currently Input and Field are the same object, but it can change in future.
type InputOptions FieldOptions
type Input Field

type ResourceInputEndpoint struct {
	Inputs []*Field `json:"inputs"`
}

func (e *ResourceInputEndpoint) findFieldIndexByName(name string) int {
	for i := 0; i < len(e.Inputs); i++ {
		if e.Inputs[i].Name == name {
			return i
		}
	}
	return -1
}

func (e *ResourceInputEndpoint) removeAtIndex(index int) *ResourceInputEndpoint {
	e.Inputs = append(e.Inputs[:index], e.Inputs[index+1:]...)
	return e
}

func (e *ResourceInputEndpoint) AddField(f Field) *ResourceInputEndpoint {
	e.Inputs = append(e.Inputs, &f)
	return e
}

func (e *ResourceInputEndpoint) RemoveField(name string) *ResourceInputEndpoint {
	idx := e.findFieldIndexByName(name)
	if idx >= 0 {
		e.removeAtIndex(idx)
	}
	return e
}

func (e *ResourceInputEndpoint) ReplaceField(name string, field Field) *ResourceInputEndpoint {
	idx := e.findFieldIndexByName(name)
	if idx >= 0 {
		e.Inputs[idx] = &field
	}
	return e
}

func NewResourceInputEndpoint() *ResourceInputEndpoint {
	return &ResourceInputEndpoint{
		Inputs: make([]*Field, 0),
	}
}
