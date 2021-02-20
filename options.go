package zepto

import (
	"github.com/go-zepto/zepto/logger"
	"github.com/go-zepto/zepto/utils"
	"github.com/go-zepto/zepto/web/renderer"
	"github.com/go-zepto/zepto/web/renderer/pongo2"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
)

type Options struct {
	Name           string
	Version        string
	Env            string
	Logger         logger.Logger
	WebpackEnabled bool
	TmplEngine     renderer.Engine
	SessionName    string
	SessionStore   sessions.Store
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	env := utils.GetEnv("ZEPTO_ENV", "development")
	opt := Options{
		Name:           "zepto",
		Version:        "latest",
		Env:            env,
		Logger:         log.New(),
		WebpackEnabled: true,
		SessionName:    "zsid",
		TmplEngine: pongo2.NewPongo2Engine(
			pongo2.TemplateDir("templates"),
			pongo2.Ext(".html"),
			pongo2.AutoReload(env == "development"),
		),
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Zepto App Name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Name of Zepto App
func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}

// Logger is the mains logger used in all zepto app
func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

func WebpackEnabled(enabled bool) Option {
	return func(o *Options) {
		o.WebpackEnabled = enabled
	}
}

// TemplateEngine  - Change the template engine implementation
func TemplateEngine(tmplEngine renderer.Engine) Option {
	return func(o *Options) {
		o.TmplEngine = tmplEngine
	}
}

// SessionName - Set the session name
func SessionName(name string) Option {
	return func(o *Options) {
		o.SessionName = name
	}
}

// SessionStore - Set the session name
func SessionStore(store sessions.Store) Option {
	return func(o *Options) {
		o.SessionStore = store
	}
}
