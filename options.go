package zepto

import "github.com/go-zepto/zepto/utils"

type Options struct {
	Name    string
	Version string
	Env     string
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	opt := Options{
		Name:    "zepto",
		Version: "latest",
		Env:     utils.GetEnv("ZEPTO_ENV", "development"),
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
