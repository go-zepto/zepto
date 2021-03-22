package gorm

import (
	"testing"

	"context"

	"github.com/go-zepto/zepto/plugins/upload/datasource"
	"github.com/go-zepto/zepto/plugins/upload/datasource/gorm/testutils"
	"github.com/go-zepto/zepto/plugins/upload/datasource/gorm/testutils/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewGormUploadDatasource(t *testing.T) {
	NewGormUploadDatasource(Options{
		DB:        testutils.SetupDB(),
		FileModel: &models.File{},
	})
}

func TestNewGormUploadDatasource_Missing_DB(t *testing.T) {
	defer func() { recover() }()
	NewGormUploadDatasource(Options{
		FileModel: &models.File{},
	})
	t.Errorf("did not panic")
}

func TestNewGormUploadDatasource_Missing_FileModel(t *testing.T) {
	defer func() { recover() }()
	NewGormUploadDatasource(Options{
		DB: testutils.SetupDB(),
	})
	t.Errorf("did not panic")
}

func TestNewGormUploadDatasource_Missing_BadFileModel_1(t *testing.T) {
	defer func() { recover() }()
	type CustomFileModel struct {
		models.Model
	}
	db := testutils.SetupDB()
	db.AutoMigrate(&CustomFileModel{})
	NewGormUploadDatasource(Options{
		DB:        testutils.SetupDB(),
		FileModel: &CustomFileModel{},
	})
	t.Errorf("did not panic")
}

func TestNewGormUploadDatasource_Missing_BadFileModel_2(t *testing.T) {
	defer func() { recover() }()
	type CustomFileModel struct {
		models.Model
		Key string `json:"key"`
	}
	db := testutils.SetupDB()
	db.AutoMigrate(&CustomFileModel{})
	NewGormUploadDatasource(Options{
		DB:        testutils.SetupDB(),
		FileModel: &CustomFileModel{},
	})
	t.Errorf("did not panic")
}

func TestNewGormUploadDatasource_Missing_BadFileModel_3(t *testing.T) {
	defer func() { recover() }()
	type CustomFileModel struct {
		models.Model
		Key string `json:"key"`
		Url string `json:"url"`
	}
	db := testutils.SetupDB()
	db.AutoMigrate(&CustomFileModel{})
	NewGormUploadDatasource(Options{
		DB:        testutils.SetupDB(),
		FileModel: &CustomFileModel{},
	})
	t.Errorf("did not panic")
}

func create(t *testing.T, ds datasource.UploadDatasource) {
	res, err := ds.Create(context.Background(), &datasource.FileData{
		Key:        "uploads/images/some-image.jpg",
		Url:        "https://localhost:8000/uploads/images/some-image.jpg",
		AccessType: "public",
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	file := res.(*models.File)
	assert.Equal(t, "uploads/images/some-image.jpg", file.Key)
	assert.Equal(t, "https://localhost:8000/uploads/images/some-image.jpg", file.Url)
	assert.Equal(t, "public", file.AccessType)
}

func TestNewGormUploadDatasource_Create(t *testing.T) {
	ds := NewGormUploadDatasource(Options{
		DB:        testutils.SetupDB(),
		FileModel: &models.File{},
	})
	create(t, ds)
}

func assertCount(t *testing.T, db *gorm.DB, key string, count int64) {
	var countRes int64
	db.Model(&models.File{}).Count(&countRes)
	assert.Equal(t, count, countRes)
}

func TestNewGormUploadDatasource_Delete(t *testing.T) {
	db := testutils.SetupDB()
	ds := NewGormUploadDatasource(Options{
		DB:        db,
		FileModel: &models.File{},
	})
	key := "uploads/images/some-image.jpg"
	assertCount(t, db, key, 0)
	create(t, ds)
	assertCount(t, db, key, 1)
	err := ds.Delete(context.Background(), key)
	assert.NoError(t, err)
	assertCount(t, db, key, 0)
}
