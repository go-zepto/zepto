package datasource

import (
	"context"

	"github.com/go-zepto/zepto/linker/filter"
)

type QueryContext struct {
	context.Context
	Filter filter.Filter
}

type SingleResult map[string]interface{}
type ListResult struct {
	Data  []map[string]interface{}
	Count int64
}

type Datasource interface {
	List(ctx QueryContext) (ListResult, error)
	Show(ctx QueryContext) (map[string]interface{}, error)
	Create(ctx QueryContext) (map[string]interface{}, error)
	Update(ctx QueryContext) (map[string]interface{}, error)
	Destroy(ctx QueryContext) (map[string]interface{}, error)
}
