package gorm

import (
	"errors"
	"reflect"

	"github.com/go-zepto/zepto/plugins/upload/datasource"
	"gorm.io/gorm"
)

var ErrNilDB = errors.New("[upload - gorm] db is nil")
var ErrNilFileModel = errors.New("[upload - gorm] file model is nil")

var FILE_MODEL_EXPECTED_COLUMNS = [...]string{
	"key",
	"url",
	"access_type",
}

type Options struct {
	DB        *gorm.DB
	FileModel interface{}
}

type GormUploadDatasource struct {
	db        *gorm.DB
	fileModel interface{}
}

func NewGormUploadDatasource(opts Options) *GormUploadDatasource {
	if opts.DB == nil {
		panic(ErrNilDB)
	}
	if opts.FileModel == nil {
		panic(ErrNilFileModel)
	}
	validateFileModel(opts.DB, opts.FileModel)
	return &GormUploadDatasource{
		db:        opts.DB,
		fileModel: opts.FileModel,
	}
}

func (d *GormUploadDatasource) createModelReflectInstance() reflect.Value {
	return reflect.New(reflect.TypeOf(d.fileModel).Elem())
}

func validateFileModel(db *gorm.DB, fileModel interface{}) {
	for _, c := range FILE_MODEL_EXPECTED_COLUMNS {
		if !db.Migrator().HasColumn(fileModel, c) {
			panic(errors.New("[upload gorm] expected " + c + " field in fileModel"))
		}
	}
}

func (d *GormUploadDatasource) Create(data *datasource.FileData) (interface{}, error) {
	obj := d.createModelReflectInstance()
	createObj := obj.Interface()
	decodeMapToStruct(data, createObj)
	query := d.db.Model(d.fileModel)
	if err := query.Create(createObj).Error; err != nil {
		return nil, err
	}
	return createObj, nil
}

func (d *GormUploadDatasource) Delete(key string) error {
	query := d.db.Model(d.fileModel)
	query = query.Where("key = ?", key)
	res := query.Delete(d.fileModel)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
