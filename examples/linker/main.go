package main

import (
	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/examples/linker/models"
	"github.com/go-zepto/zepto/plugins/linker"
	gormds "github.com/go-zepto/zepto/plugins/linker/datasource/gorm"
	"github.com/go-zepto/zepto/plugins/linker/filter"
	"github.com/go-zepto/zepto/web"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(
		&models.Author{},
		&models.Book{},
	)
	return db
}

func main() {
	db := SetupDB()
	z := zepto.NewZepto(
		zepto.Name("books-api"),
		zepto.Version("1.0.0"),
	)

	z.AddPlugin(linker.NewLinkerPlugin(linker.Options{
		Path: "/api",
		Resources: []linker.Resource{
			{
				Name:       "Author",
				Datasource: gormds.NewGormDatasource(db, &models.Author{}),
			},
			{
				Name:       "Book",
				Datasource: gormds.NewGormDatasource(db, &models.Book{}),
			},
		},
	}))

	z.Get("/first-author", func(ctx web.Context) error {
		l := linker.InstanceFromCtx(ctx)
		var author *models.Author
		l.RepositoryDecoder("Author").FindOne(ctx, &filter.Filter{}, &author)
		return ctx.RenderJson(map[string]*models.Author{
			"first_author": author,
		})
	})

	z.SetupHTTP("0.0.0.0:8000")
	z.Start()
}
