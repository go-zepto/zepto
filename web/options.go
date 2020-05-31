package web

import (
	"github.com/go-zepto/zepto/broker"
	"github.com/go-zepto/zepto/web/renderer"
	log "github.com/sirupsen/logrus"
)

type Options struct {
	broker         *broker.Broker
	logger         *log.Logger
	env            string
	webpackEnabled bool
	tmplEngine     renderer.Engine
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	opt := Options{
		webpackEnabled: true,
		env:            "development",
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Zepto Broker
func Broker(b *broker.Broker) Option {
	return func(o *Options) {
		o.broker = b
	}
}

// Zepto Logger
func Logger(l *log.Logger) Option {
	return func(o *Options) {
		o.logger = l
	}
}

// Zepto Env
func Env(e string) Option {
	return func(o *Options) {
		o.env = e
	}
}

func WebpackEnabled(enabled bool) Option {
	return func(o *Options) {
		o.webpackEnabled = enabled
	}
}

// TemplateEngine  - Change the template engine implementation
func TemplateEngine(tmplEngine renderer.Engine) Option {
	return func(o *Options) {
		o.tmplEngine = tmplEngine
	}
}
