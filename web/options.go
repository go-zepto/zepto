package web

import (
	"github.com/go-zepto/zepto/broker"
	log "github.com/sirupsen/logrus"
)

type Options struct {
	broker         *broker.Broker
	logger         *log.Logger
	env            string
	webpackEnabled bool
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

// Zepto Logger
func WebpackEnabled(enabled bool) Option {
	return func(o *Options) {
		o.webpackEnabled = enabled
	}
}
