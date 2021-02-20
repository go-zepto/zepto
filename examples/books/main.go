package main

import (
	"log"
	"os"
	"time"

	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/examples/books/models"
	"github.com/go-zepto/zepto/linker"
	gormds "github.com/go-zepto/zepto/linker/datasource/gorm"
	"github.com/go-zepto/zepto/web"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDB() *gorm.DB {
	l := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: l,
	})
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
		zepto.Version("0.0.1"),
	)

	app := z.NewWeb()

	api := app.Router("/api", web.Hosts("localhost:8000"))

	lr := linker.NewLinker(api)

	lr.AddResource(linker.Resource{
		Name:       "Author",
		Datasource: gormds.NewGormDatasource(db, &models.Author{}),
	})

	lr.AddResource(linker.Resource{
		Name:       "Book",
		Datasource: gormds.NewGormDatasource(db, &models.Book{}),
	})

	z.SetupHTTP("0.0.0.0:8000", app)
	z.Start()
}
