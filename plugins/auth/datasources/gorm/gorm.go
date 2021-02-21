package gormauth

import (
	"github.com/go-zepto/zepto/plugins/auth"
	"gorm.io/gorm"
)

type GormAuthDatasource struct {
	db    *gorm.DB
	model interface{}
}

type GormAuthDatasourceOptions struct {
	DB    *gorm.DB
	Model interface{}
}

func NewGormAuthDatasoruce(opts GormAuthDatasourceOptions) *GormAuthDatasource {
	gads := GormAuthDatasource{
		db:    opts.DB,
		model: opts.Model,
	}
	return &gads
}

func (gads *GormAuthDatasource) Auth(username string, password string) (auth.PID, error) {
	return auth.PID(100), nil
}
