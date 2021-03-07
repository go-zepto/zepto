package linker

import (
	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/web"
)

type Options struct {
	Path        string
	Hosts       []string
	SetupRouter func(router *web.Router)
	Resources   []Resource
}

func NewLinkerPlugin(opts Options) *LinkerPlugin {
	if opts.Path == "" {
		opts.Path = "/api"
	}
	return &LinkerPlugin{
		opts: opts,
	}
}

type LinkerPlugin struct {
	opts   Options
	linker *Linker
	router *web.Router
}

func (l *LinkerPlugin) Name() string {
	return "linker"
}

func (l *LinkerPlugin) Instance() interface{} {
	return LinkerInstance(l.linker)
}

func (l *LinkerPlugin) PrependMiddlewares() []web.MiddlewareFunc {
	return nil
}

func (l *LinkerPlugin) AppendMiddlewares() []web.MiddlewareFunc {
	return nil
}

func (l *LinkerPlugin) OnCreated(z *zepto.Zepto) {
	l.router = z.Router(l.opts.Path, web.Hosts(l.opts.Hosts...))
	if l.opts.SetupRouter != nil {
		l.opts.SetupRouter(l.router)
	}
	l.linker = NewLinker(l.router)
}

func (l *LinkerPlugin) OnZeptoInit(z *zepto.Zepto) {
	for _, r := range l.opts.Resources {
		l.linker.AddResource(r)
		z.Logger().Debugf("[linker] Resource %s added", r.Name)
	}
}

func (l *LinkerPlugin) OnZeptoStart(z *zepto.Zepto) {}

func (l *LinkerPlugin) OnZeptoStop(z *zepto.Zepto) {}
