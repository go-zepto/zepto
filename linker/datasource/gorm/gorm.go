package gorm

import (
	"reflect"

	"github.com/go-zepto/zepto/linker/datasource"
	"github.com/go-zepto/zepto/linker/filter/where"
	lutils "github.com/go-zepto/zepto/linker/utils"
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
	if err := query.First(&dest).Error; err != nil {
		return nil, err
	}
	return &dest, nil
}

func (g *GormDatasource) createModelReflectInstance() reflect.Value {
	return reflect.New(reflect.TypeOf(g.Model).Elem())
}

func (g *GormDatasource) Create(ctx datasource.QueryContext, data map[string]interface{}) (*map[string]interface{}, error) {
	obj := g.createModelReflectInstance()
	createObj := obj.Interface()
	lutils.DecodeMapToStruct(data, createObj)
	query := g.DB.Model(g.Model)
	if err := query.Create(createObj).Error; err != nil {
		return nil, err
	}
	result := lutils.DecodeStructToMap(createObj)
	return &result, nil
}

func (g *GormDatasource) Update(ctx datasource.QueryContext, data map[string]interface{}) (datasource.ManyAffectedResult, error) {
	query := g.DB.Model(g.Model)
	query, err := g.ApplyWhere(ctx, query)
	if err != nil {
		return datasource.ManyAffectedResult{}, err
	}
	res := query.Updates(data)
	if res.Error != nil {
		return datasource.ManyAffectedResult{}, res.Error
	}
	return datasource.ManyAffectedResult{
		TotalAffected: res.RowsAffected,
	}, nil
}

func (g *GormDatasource) Destroy(ctx datasource.QueryContext) (datasource.ManyAffectedResult, error) {
	query := g.DB.Model(g.Model)
	query, err := g.ApplyWhere(ctx, query)
	if err != nil {
		return datasource.ManyAffectedResult{}, err
	}
	res := query.Delete(g.Model)
	if res.Error != nil {
		return datasource.ManyAffectedResult{}, res.Error
	}
	return datasource.ManyAffectedResult{
		TotalAffected: res.RowsAffected,
	}, nil
}
