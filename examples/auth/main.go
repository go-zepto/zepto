package main

import (
	"os"

	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/examples/auth/models"
	"github.com/go-zepto/zepto/plugins/auth"
	gormds "github.com/go-zepto/zepto/plugins/auth/datasources/gorm"
	"github.com/go-zepto/zepto/plugins/mailer"
	"github.com/go-zepto/zepto/plugins/mailer/providers/sendgrid"
	"github.com/go-zepto/zepto/web"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, email string, username string, password string) {
	pwd, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	u := models.User{
		Email:        email,
		Username:     username,
		PasswordHash: string(pwd),
	}
	db.Create(&u)
}

func SetupDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(
		&models.User{},
	)
	CreateUser(db, "clark.kent@gozepto.com", "clark.kent", "i.am.superman")
	CreateUser(db, "bruce.wayne@gozepto.com", "bruce.wayne", "i.am.batman")
	return db
}

func main() {
	db := SetupDB()
	z := zepto.NewZepto(
		zepto.Name("auth-api"),
		zepto.Version("0.0.1"),
	)

	z.AddPlugin(mailer.NewMailerPlugin(mailer.Options{
		Mailer: sendgrid.NewSendgridProvider(sendgrid.Settings{
			ApiKey: os.Getenv("SENDGRID_API_TOKEN"),
			DefaultFrom: &mailer.Email{
				Name:    "Go Zepto",
				Address: "no-reply@gozepto.com",
			},
		}),
	}))

	z.AddPlugin(auth.NewAuthTokenPlugin(auth.AuthTokenOptions{
		Datasource: gormds.NewGormAuthDatasoruce(gormds.GormAuthDatasourceOptions{
			DB:        db,
			UserModel: &models.User{},
		}),
	}))

	z.Get("/me", func(ctx web.Context) error {
		a := auth.InstanceFromCtx(ctx)
		pid := a.LoggedPIDFromCtx(ctx)
		if pid == nil {
			ctx.SetStatus(401)
			return ctx.RenderJson(map[string]string{
				"error": "unauthorized",
			})
		}
		var me models.User
		err := db.First(&me, pid).Error
		if err != nil {
			ctx.SetStatus(500)
			return ctx.RenderJson(map[string]string{
				"error": "internal server error",
			})
		}
		return ctx.RenderJson(me)
	})

	z.SetupHTTP("0.0.0.0:8000")
	z.Start()
}
