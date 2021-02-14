package repository

import (
	"context"

	"github.com/go-zepto/zepto/linker/datasource"
	"github.com/go-zepto/zepto/linker/filter"
)

type Repository struct {
	ds datasource.Datasource
}

func (r *Repository) FindById(ctx context.Context, id interface{}) (*map[string]interface{}, error) {
	where := map[string]interface{}{
		"id": map[string]interface{}{
			"eq": id,
		},
	}
	return r.ds.FindOne(datasource.QueryContext{
		Context: ctx,
		Filter: &filter.Filter{
			Where: &where,
		},
	})
}
