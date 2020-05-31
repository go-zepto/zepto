package pongo2

type Options struct {
	templateDir string
	autoReload  bool
	ext         string
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	opt := Options{
		templateDir: "templates",
		ext:         ".html",
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// TemplateDir - Template directory. Default=templates
func TemplateDir(templateDir string) Option {
	return func(o *Options) {
		o.templateDir = templateDir
	}
}

// AutoReload - Template directory.
func AutoReload(autoReload bool) Option {
	return func(o *Options) {
		o.autoReload = autoReload
	}
}

// Ext - Template extension. Default=.html
func Ext(ext string) Option {
	return func(o *Options) {
		o.ext = ext
	}
}
