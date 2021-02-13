package gorm

import (
	"github.com/go-zepto/zepto/linker/datasource"
	"github.com/go-zepto/zepto/linker/filter/where"
	"gorm.io/gorm"
)

type GormDatasource struct {
	DB    *gorm.DB
	Model interface{}
}

func NewGormDatasource(db *gorm.DB, model interface{}) *GormDatasource {
	return &GormDatasource{
		DB:    db,
		Model: model,
	}
}

func (g *GormDatasource) List(ctx datasource.QueryContext) (datasource.ListResult, error) {
	dest := []map[string]interface{}{}
	var count int64
	query := g.DB.Model(g.Model)
	where := where.NewFromMap(ctx.Filter.Where)
	sqlWhere, err := where.ToSQL()
	if err != nil {
		return datasource.ListResult{}, err
	}
	query = query.Where(sqlWhere.Text, sqlWhere.Vars...)
	countQuery := query
	query = query.Limit(10).Offset(0)
	if err := query.Find(&dest).Error; err != nil {
		return datasource.ListResult{}, err
	}
	if err := countQuery.Count(&count).Error; err != nil {
		return datasource.ListResult{}, err
	}
	return datasource.ListResult{
		Data:  dest,
		Count: count,
	}, nil
}

func (g *GormDatasource) Show(ctx datasource.QueryContext) (map[string]interface{}, error) {
	panic("not implemented")
}

func (g *GormDatasource) Create(ctx datasource.QueryContext) (map[string]interface{}, error) {
	panic("not implemented")
}

func (g *GormDatasource) Update(ctx datasource.QueryContext) (map[string]interface{}, error) {
	panic("not implemented")
}

func (g *GormDatasource) Destroy(ctx datasource.QueryContext) (map[string]interface{}, error) {
	panic("not implemented")
}
