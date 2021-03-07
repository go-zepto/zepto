package utils

import (
	"sync"

	"github.com/go-zepto/zepto/web/renderer"
)

func RenderTemplateToString(engine renderer.Engine, template string, vars map[string]interface{}) (string, error) {
	w := MailerWriter{}
	varsMap := sync.Map{}
	for key, v := range vars {
		varsMap.Store(key, v)
	}
	err := engine.Render(&w, 200, template, &varsMap)
	if err != nil {
		return "", err
	}
	html := string(w.Value())
	return html, nil
}
