package main

import (
	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/web"
)

func main() {
	z := zepto.NewZepto()

	z.Get("/hello", func(ctx web.Context) error {
		ctx.SetStatus(200)
		return ctx.RenderJson(map[string]interface{}{
			"hello": "world",
		})
	})

	z.Start()
}
