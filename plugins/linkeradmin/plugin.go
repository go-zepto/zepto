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

// Currently Input and Field are the same object, but it can change in future.
type InputOptions FieldOptions
type Input Field

type LinkerResource struct {
	Name         string  `json:"name"`
	Endpoint     string  `json:"endpoint"`
	ListFields   []Field `json:"list_fields"`
	CreateInputs []Input `json:"create_inputs"`
	UpdateInputs []Input `json:"update_inputs"`
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
	res := make([]LinkerResource, 0)
	for _, r := range opts.LinkerResources {
		if r.ListFields == nil {
			r.ListFields = make([]Field, 0)
		}
		if r.CreateInputs == nil {
			r.CreateInputs = make([]Input, 0)
		}
		if r.UpdateInputs == nil {
			r.UpdateInputs = make([]Input, 0)
		}
		res = append(res, r)
	}
	return &LinkerAdminPlugin{
		Schema: &Schema{
			Resources: res,
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
