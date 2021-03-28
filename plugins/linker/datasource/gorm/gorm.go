package gorm

import (
	"reflect"
	"strings"

	"github.com/go-zepto/zepto/plugins/linker/datasource"
	"github.com/go-zepto/zepto/plugins/linker/filter/where"
	lutils "github.com/go-zepto/zepto/plugins/linker/utils"
	"gorm.io/gorm"
)

type GormDatasource struct {
	DB         *gorm.DB
	Model      interface{}
	Properties datasource.Properties
	fields     map[string]datasource.Field
}

func getFields(db *gorm.DB, model interface{}) (*map[string]datasource.Field, error) {
	types, err := db.Migrator().ColumnTypes(model)
	if err != nil {
		return nil, err
	}
	fields := make(map[string]datasource.Field)
	for _, t := range types {
		nullable, _ := t.Nullable()
		fields[t.Name()] = datasource.Field{
			Name:     t.Name(),
			Type:     strings.ToLower(t.DatabaseTypeName()),
			Nullable: nullable,
		}
	}
	return &fields, nil
}

func NewGormDatasource(db *gorm.DB, model interface{}) *GormDatasource {
	fields, err := getFields(db, model)
	if err != nil {
		panic(err)
	}
	return &GormDatasource{
		DB:    db,
		Model: model,
		Properties: datasource.Properties{
			Skip:  0,
			Limit: 10,
		},
		fields: *fields,
	}
}

func (g *GormDatasource) getFieldList() []string {
	fieldList := make([]string, 0)
	for fieldName := range g.fields {
		fieldList = append(fieldList, fieldName)
	}
	return fieldList
}

func (g *GormDatasource) ApplyWhere(ctx datasource.QueryContext, query *gorm.DB) (*gorm.DB, error) {
	if ctx.Filter != nil && ctx.Filter.Where != nil {
		allowedFields := g.getFieldList()
		where := where.NewFromMapWithAllowedFields(*ctx.Filter.Where, allowedFields)
		sqlWhere, err := where.ToSQL()
		if err != nil {
			return nil, err
		}
		query = query.Where(sqlWhere.Text, sqlWhere.Vars...)
	}
	return query, nil
}

func (g *GormDatasource) ApplyInclude(ctx datasource.QueryContext, query *gorm.DB) (*gorm.DB, error) {
	if ctx.Filter != nil && ctx.Filter.Include != nil {
		for _, include := range ctx.Filter.Include {
			if include.Where == nil {
				query = query.Preload(include.Relation)
			} else {
				where := where.NewFromMap(*include.Where)
				sqlWhere, err := where.ToSQL()
				if err != nil {
					return nil, err
				}
				args := []interface{}{}
				args = append(args, sqlWhere.Text)
				for _, v := range sqlWhere.Vars {
					args = append(args, v)
				}
				query = query.Preload(include.Relation, args...)
			}
		}
	}
	return query, nil
}

func (g *GormDatasource) createModelReflectInstance() reflect.Value {
	return reflect.New(reflect.TypeOf(g.Model).Elem())
}

func (g *GormDatasource) createModelReflectInstanceSlice() reflect.Value {
	elemType := reflect.TypeOf(g.Model)
	return reflect.New(reflect.SliceOf(elemType))
}

func (g *GormDatasource) Fields() map[string]datasource.Field {
	return g.fields
}

func (g *GormDatasource) Find(ctx datasource.QueryContext) (*datasource.ListResult, error) {
	obj := g.createModelReflectInstanceSlice()
	dest := obj.Interface()
	var count int64
	query := g.DB.Model(g.Model)
	query, err := g.ApplyInclude(ctx, query)
	if err != nil {
		return nil, err
	}
	query, err = g.ApplyWhere(ctx, query)
	if err != nil {
		return nil, err
	}
	countQuery := query
	if err := countQuery.Count(&count).Error; err != nil {
		return nil, err
	}
	query = query.Limit(int(g.Properties.GetLimit(ctx))).Offset(int(g.Properties.GetSkip(ctx)))
	if err := query.Find(dest).Error; err != nil {
		return nil, err
	}
	result := lutils.DecodeStructToMapList(dest)
	return &datasource.ListResult{
		Data:  result,
		Count: count,
	}, nil
}

func (g *GormDatasource) FindOne(ctx datasource.QueryContext) (*map[string]interface{}, error) {
	obj := g.createModelReflectInstance()
	dest := obj.Interface()
	query := g.DB.Model(g.Model)
	query, err := g.ApplyInclude(ctx, query)
	if err != nil {
		return nil, err
	}
	query, err = g.ApplyWhere(ctx, query)
	if err != nil {
		return nil, err
	}
	if err := query.First(dest).Error; err != nil {
		return nil, err
	}
	result := lutils.DecodeStructToMap(dest)
	return &result, nil
}

func (g *GormDatasource) Create(ctx datasource.QueryContext, data interface{}) (*map[string]interface{}, error) {
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

func (g *GormDatasource) Update(ctx datasource.QueryContext, data interface{}) (datasource.ManyAffectedResult, error) {
	obj := g.createModelReflectInstance()
	updateObj := obj.Interface()
	lutils.DecodeMapToStruct(data, updateObj)
	query := g.DB.Model(g.Model)
	query, err := g.ApplyWhere(ctx, query)
	if err != nil {
		return datasource.ManyAffectedResult{}, err
	}
	res := query.Updates(updateObj)
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
