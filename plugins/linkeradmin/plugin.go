package linkeradmin

import (
	"embed"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/web"
)

//go:embed webapp/build/*
var webappBuild embed.FS

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
	Path            string
}

type Schema struct {
	Resources []LinkerResource `json:"resources"`
}

type LinkerAdminPlugin struct {
	Schema *Schema
	path   string
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
	if opts.Path == "" {
		opts.Path = "/admin"
	}
	return &LinkerAdminPlugin{
		Schema: &Schema{
			Resources: res,
		},
		path: opts.Path,
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
	l.router = z.Router(l.path)
}

func (l *LinkerAdminPlugin) OnZeptoInit(z *zepto.Zepto) {
	l.router.Get("/_schema", func(ctx web.Context) error {
		return ctx.RenderJson(l.Schema)
	})

	l.router.Get("/", func(ctx web.Context) error {
		fmt.Println(ctx.Request().URL.Path, l.path)
		// Root index.html
		indexHTMLBytes, err := webappBuild.ReadFile("webapp/build/index.html")
		if err != nil {
			return errors.New("could not load admin")
		}
		indexHTML := string(indexHTMLBytes)
		fmt.Println(indexHTML)
		indexHTML = strings.ReplaceAll(indexHTML, "/admin", l.path)
		fmt.Println(indexHTML)
		ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err = fmt.Fprint(ctx.Response(), indexHTML)
		return err
	})

	l.router.Get("/{rest:.*}", func(ctx web.Context) error {
		if ctx.Request().URL.Path == l.path+"/" {
			return ctx.Redirect(l.path)
		}
		filePath := strings.Replace(ctx.Request().URL.Path, l.path, "", 1)
		req, _ := http.NewRequest("GET", "webapp/build"+filePath, nil)
		handler := http.FileServer(http.FS(webappBuild))
		handler.ServeHTTP(ctx.Response(), req)
		return nil
	})
}

func (l *LinkerAdminPlugin) OnZeptoStart(z *zepto.Zepto) {}

func (l *LinkerAdminPlugin) OnZeptoStop(z *zepto.Zepto) {}
