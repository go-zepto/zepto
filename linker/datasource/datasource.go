package datasource

import (
	"context"

	"github.com/go-zepto/zepto/linker/filter"
)

type QueryContext struct {
	context.Context
	Filter *filter.Filter
}

type SingleResult map[string]interface{}
type ListResult struct {
	Data  []map[string]interface{}
	Count int64
}

type ManyAffectedResult struct {
	TotalAffected int64
}

type Datasource interface {
	Find(ctx QueryContext) (*ListResult, error)
	FindOne(ctx QueryContext) (*map[string]interface{}, error)
	Create(ctx QueryContext, data map[string]interface{}) (*map[string]interface{}, error)
	Update(ctx QueryContext, data map[string]interface{}) (ManyAffectedResult, error)
	Destroy(ctx QueryContext) (ManyAffectedResult, error)
}

type Properties struct {
	Skip  int64
	Limit int64
}

func (d *Properties) GetSkip(ctx QueryContext) int64 {
	f := ctx.Filter
	if f != nil && f.Skip != nil {
		return *f.Skip
	}
	return d.Skip
}

func (d *Properties) GetLimit(ctx QueryContext) int64 {
	f := ctx.Filter
	if f != nil && f.Limit != nil {
		return *f.Limit
	}
	return d.Limit
}
