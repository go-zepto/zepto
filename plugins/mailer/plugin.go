package mailer

import (
	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/web"
	"github.com/go-zepto/zepto/web/renderer"
)

type Options struct {
	Mailer Mailer
}

func NewMailerPlugin(opts Options) *MailerPlugin {
	return &MailerPlugin{
		mailer: opts.Mailer,
	}
}

type MailerPlugin struct {
	mailer   Mailer
	instance MailerInstance
}

func (mp *MailerPlugin) Name() string {
	return "mailer"
}

func (mp *MailerPlugin) Instance() interface{} {
	return MailerInstance(mp.mailer)
}

func (mp *MailerPlugin) PrependMiddlewares() []web.MiddlewareFunc {
	return nil
}

func (mp *MailerPlugin) AppendMiddlewares() []web.MiddlewareFunc {
	return nil
}

func (mp *MailerPlugin) OnCreated(z *zepto.Zepto) {}

func (mp *MailerPlugin) OnZeptoInit(z *zepto.Zepto) {
	var engine renderer.Engine
	if z.App != nil {
		engine = z.App.RendererEngine()
	}
	mp.mailer.Init(&InitOptions{
		RendererEngine: engine,
	})
}

func (mp *MailerPlugin) OnZeptoStart(z *zepto.Zepto) {}

func (mp *MailerPlugin) OnZeptoStop(z *zepto.Zepto) {}
