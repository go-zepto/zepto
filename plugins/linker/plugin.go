package linker

import (
	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/web"
)

type Options struct {
	Linker *Linker
}

func NewLinkerPlugin(opts Options) *LinkerPlugin {
	return &LinkerPlugin{
		opts: opts,
	}
}

type LinkerPlugin struct {
	opts   Options
	linker *Linker
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
	l.linker = l.opts.Linker
}

func (l *LinkerPlugin) OnZeptoInit(z *zepto.Zepto) {}

func (l *LinkerPlugin) OnZeptoStart(z *zepto.Zepto) {}

func (l *LinkerPlugin) OnZeptoStop(z *zepto.Zepto) {}
