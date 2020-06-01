package web

import (
	"github.com/go-zepto/zepto/broker"
	"github.com/go-zepto/zepto/logger"
	"github.com/go-zepto/zepto/utils"
	"github.com/go-zepto/zepto/web/renderer"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
)

type Options struct {
	broker         *broker.Broker
	logger         logger.Logger
	env            string
	webpackEnabled bool
	tmplEngine     renderer.Engine
	sessionName    string
	sessionStore   sessions.Store
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	opt := Options{
		webpackEnabled: true,
		env:            utils.GetEnv("ZEPTO_ENV", "development"),
		sessionName:    "zsid",
		logger:         log.New(),
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
func Logger(l logger.Logger) Option {
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

// SessionName - Set the session name
func SessionName(name string) Option {
	return func(o *Options) {
		o.sessionName = name
	}
}

// SessionStore - Set the session name
func SessionStore(store sessions.Store) Option {
	return func(o *Options) {
		o.sessionStore = store
	}
}
