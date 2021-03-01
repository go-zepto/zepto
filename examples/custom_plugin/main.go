package main

import (
	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/examples/custom_plugin/plugins/counter"
	"github.com/go-zepto/zepto/web"
)

func main() {
	z := zepto.NewZepto(zepto.Name("app-custom-plugin"))
	z.AddPlugin(counter.NewCounterPlugin())
	z.SetupHTTP("0.0.0.0:8000")

	// You can access the plugin instance using the Controller Context
	z.Get("/", func(ctx web.Context) error {
		counter := (ctx.PluginInstance("counter")).(*counter.Counter)
		counter.Add(100)
		return ctx.RenderJson(map[string]int{
			"counter": counter.Value(),
		})
	})

	z.Start()
}
