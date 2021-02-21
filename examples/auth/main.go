package main

import (
	"log"
	"os"
	"time"

	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/examples/auth/models"
	"github.com/go-zepto/zepto/plugins/auth"
	gormauth "github.com/go-zepto/zepto/plugins/auth/datasources/gorm"
	inmemory "github.com/go-zepto/zepto/plugins/auth/stores/inmemory"
	"github.com/go-zepto/zepto/plugins/auth/token_encoders/uuid"
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
		&models.User{},
	)
	return db
}

func main() {
	db := SetupDB()
	z := zepto.NewZepto(
		zepto.Name("books-api"),
		zepto.Version("0.0.1"),
	)

	authToken := auth.NewAuthToken(auth.AuthTokenOptions{
		Datasource: gormauth.NewGormAuthDatasoruce(gormauth.GormAuthDatasourceOptions{
			DB:    db,
			Model: &models.User{},
		}),
		TokenEncoder: uuid.NewUUIDTokenEncoder(),
		Store:        inmemory.NewInMemoryStore(),
	})

	authToken.Create(z)

	z.Get("/me", func(ctx web.Context) error {
		pid := ctx.Value("auth_user_pid")
		if pid == nil {
			ctx.SetStatus(401)
			return ctx.RenderJson(map[string]interface{}{
				"error": "not authorized",
			})
		}
		userId := pid.(int)
		return ctx.RenderJson(map[string]interface{}{
			"user_id": userId,
		})
	})

	authToken.Init(z)

	z.SetupHTTP("0.0.0.0:8000")
	z.Start()
}
