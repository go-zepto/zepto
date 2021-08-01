package main

import (
	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/web"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	z := zepto.NewZepto()
	z.DB().AutoMigrate(&Product{})
	// z.DB().Create(&Product{Model: gorm.Model{ID: 1}, Code: "D42", Price: 100})
	// z.DB().Create(&Product{Model: gorm.Model{ID: 1}, Code: "D42", Price: 100})

	z.Get("/hello", func(ctx web.Context) error {
		ctx.SetStatus(200)
		return ctx.RenderJson(map[string]interface{}{
			"hello": "world",
		})
	})

	z.Start()
}
