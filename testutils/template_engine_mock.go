package testutils

import (
	"net/http"
	"sync"
)

type EngineMockRenderCallArgs struct {
	status   int
	template string
	data     *sync.Map
}

type EngineMock struct {
	InitCalled     bool
	RenderCalled   bool
	RenderCallArgs *EngineMockRenderCallArgs
}

func (e *EngineMock) Init() error {
	e.InitCalled = true
	return nil
}

func (e *EngineMock) Render(w http.ResponseWriter, status int, template string, data *sync.Map) error {
	e.RenderCalled = true
	e.RenderCallArgs = &EngineMockRenderCallArgs{
		status:   status,
		template: template,
		data:     data,
	}
	w.WriteHeader(200)
	w.Write([]byte("Mocked Template!"))
	return nil
}
