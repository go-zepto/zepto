package renderer

import (
	"net/http"
	"sync"
)

type EngineOptions struct {
	// TemplateDir is where all templates is located. Default=templates
	TemplateDir string
	// Default file extension for templates. Default=.html
	Ext string
	// AutoReload if enabled, reload the template from file on each change.
	AutoReload bool
}

type Engine interface {
	// Init should be called for prepare the template engine. Load the template dir, for example.
	Init() error
	// Render is used to render a template
	Render(w http.ResponseWriter, status int, template string, data *sync.Map) error
}
