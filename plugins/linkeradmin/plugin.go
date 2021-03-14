package linkeradmin

import (
	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/web"
)

type FieldOptions = map[string]interface{}

type Field struct {
	Name    string       `json:"name"`
	Type    string       `json:"type"`
	Options FieldOptions `json:"options"`
}

type LinkerResource struct {
	Name     string  `json:"name"`
	Endpoint string  `json:"endpoint"`
	Fields   []Field `json:"fields"`
}

type Options struct {
	LinkerResources []LinkerResource
}

type Schema struct {
	Resources []LinkerResource `json:"resources"`
}

type LinkerAdminPlugin struct {
	Schema *Schema
	router *web.Router
}

func NewLinkerAdminPlugin(opts Options) *LinkerAdminPlugin {
	return &LinkerAdminPlugin{
		Schema: &Schema{
			Resources: opts.LinkerResources,
		},
	}
}

func (l *LinkerAdminPlugin) Name() string {
	return "linkeradmin"
}

func (l *LinkerAdminPlugin) Instance() interface{} {
	return nil
}

func (l *LinkerAdminPlugin) PrependMiddlewares() []web.MiddlewareFunc {
	return nil
}

func (l *LinkerAdminPlugin) AppendMiddlewares() []web.MiddlewareFunc {
	return nil
}

func (l *LinkerAdminPlugin) OnCreated(z *zepto.Zepto) {
	l.router = z.Router("/admin")
}

func (l *LinkerAdminPlugin) OnZeptoInit(z *zepto.Zepto) {
	l.router.Get("/_schema", func(ctx web.Context) error {
		return ctx.RenderJson(l.Schema)
	})
}

func (l *LinkerAdminPlugin) OnZeptoStart(z *zepto.Zepto) {}

func (l *LinkerAdminPlugin) OnZeptoStop(z *zepto.Zepto) {}
