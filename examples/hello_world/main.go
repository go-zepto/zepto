package main

import (
	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/web"
)

func main() {
	z := zepto.NewZepto(
		zepto.Name("hello-world-app"),
		zepto.Version("0.0.1"),
	)

	api := z.Router("/")

	api.Get("/", func(ctx web.Context) error {
		panic("Opssss")
		ctx.SetStatus(400)
		return ctx.RenderJson(map[string]string{"error": "bad request"})
	})

	z.SetupHTTP("0.0.0.0:8000")
	z.Start()
}
