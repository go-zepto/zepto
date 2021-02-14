package gorm

import (
	"reflect"

	"github.com/go-zepto/zepto/linker/datasource"
	"github.com/go-zepto/zepto/linker/datasource/gorm/testutils"
	"github.com/go-zepto/zepto/linker/datasource/gorm/utils"
	"github.com/go-zepto/zepto/linker/filter/where"
	"gorm.io/gorm"
)

type GormDatasource struct {
	DB         *gorm.DB
	Model      interface{}
	Properties datasource.Properties
}

func NewGormDatasource(db *gorm.DB, model interface{}) *GormDatasource {
	return &GormDatasource{
		DB:    db,
		Model: model,
		Properties: datasource.Properties{
			Skip:  0,
			Limit: 10,
		},
	}
}

func (g *GormDatasource) ApplyWhere(ctx datasource.QueryContext, query *gorm.DB) (*gorm.DB, error) {
	if ctx.Filter != nil && ctx.Filter.Where != nil {
		where := where.NewFromMap(*ctx.Filter.Where)
		sqlWhere, err := where.ToSQL()
		if err != nil {
			return nil, err
		}
		query = query.Where(sqlWhere.Text, sqlWhere.Vars...)
	}
	return query, nil
}

func (g *GormDatasource) Find(ctx datasource.QueryContext) (*datasource.ListResult, error) {
	dest := []map[string]interface{}{}
	var count int64
	query := g.DB.Model(g.Model)
	query, err := g.ApplyWhere(ctx, query)
	if err != nil {
		return nil, err
	}
	countQuery := query
	if err := countQuery.Count(&count).Error; err != nil {
		return nil, err
	}
	query = query.Limit(int(g.Properties.GetLimit(ctx))).Offset(int(g.Properties.GetSkip(ctx)))
	if err := query.Find(&dest).Error; err != nil {
		return nil, err
	}
	return &datasource.ListResult{
		Data:  dest,
		Count: count,
	}, nil
}

func (g *GormDatasource) FindOne(ctx datasource.QueryContext) (*map[string]interface{}, error) {
	dest := map[string]interface{}{}
	query := g.DB.Model(g.Model)
	query, err := g.ApplyWhere(ctx, query)
	if err != nil {
		return nil, err
	}
	if err := query.Find(&dest).Error; err != nil {
		return nil, err
	}
	return &dest, nil
}

func (g *GormDatasource) itemGen1() interface{} {
	// d := time.Date(1990, 5, 5, 0, 0, 0, 0, time.Local)
	return &testutils.Person{}
}

func (g *GormDatasource) createModelReflectInstance() reflect.Value {
	return reflect.New(reflect.TypeOf(g.Model).Elem())
}

func (g *GormDatasource) Create(ctx datasource.QueryContext, data map[string]interface{}) (*map[string]interface{}, error) {
	obj := g.createModelReflectInstance()
	createObj := obj.Interface()
	utils.DecodeMapToStruct(data, createObj)
	query := g.DB.Model(g.Model)
	if err := query.Create(createObj).Error; err != nil {
		return nil, err
	}
	result := utils.DecodeStructToMap(createObj)
	return &result, nil
}

func (g *GormDatasource) Update(ctx datasource.QueryContext, data map[string]interface{}) (*map[string]interface{}, error) {
	panic("not implemented")
}

func (g *GormDatasource) Destroy(ctx datasource.QueryContext) (*map[string]interface{}, error) {
	panic("not implemented")
}
