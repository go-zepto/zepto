package linkeradmin

import "github.com/go-zepto/zepto/plugins/linkeradmin/fields"

type ResourceInputEndpoint struct {
	Inputs []*fields.Input `json:"inputs"`
}

func (e *ResourceInputEndpoint) findInputIndexByName(name string) int {
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

func (e *ResourceInputEndpoint) AddInput(i fields.Input) *ResourceInputEndpoint {
	e.Inputs = append(e.Inputs, &i)
	return e
}

func (e *ResourceInputEndpoint) RemoveInput(name string) *ResourceInputEndpoint {
	idx := e.findInputIndexByName(name)
	if idx >= 0 {
		e.removeAtIndex(idx)
	}
	return e
}

func (e *ResourceInputEndpoint) ReplaceInput(name string, input fields.Input) *ResourceInputEndpoint {
	idx := e.findInputIndexByName(name)
	if idx >= 0 {
		e.Inputs[idx] = &input
	}
	return e
}

func NewResourceInputEndpoint() *ResourceInputEndpoint {
	return &ResourceInputEndpoint{
		Inputs: make([]*fields.Input, 0),
	}
}
