package repository

import (
	"context"

	"github.com/go-zepto/zepto/linker/datasource"
	"github.com/go-zepto/zepto/linker/filter"
)

type Repository struct {
	ds datasource.Datasource
}

type RepositoryConfig struct {
	Datasource datasource.Datasource
}

func NewRepository(config RepositoryConfig) *Repository {
	return &Repository{
		ds: config.Datasource,
	}
}

func (r *Repository) FindById(ctx context.Context, id interface{}) (*SingleResult, error) {
	where := map[string]interface{}{
		"id": map[string]interface{}{
			"eq": id,
		},
	}
	res, err := r.ds.FindOne(datasource.QueryContext{
		Context: ctx,
		Filter: &filter.Filter{
			Where: &where,
		},
	})
	if err != nil {
		return nil, err
	}
	rr := SingleResult(*res)
	return &rr, err
}

func (r *Repository) FindOne(ctx context.Context, filter *filter.Filter) (*SingleResult, error) {
	res, err := r.ds.FindOne(datasource.QueryContext{
		Context: ctx,
		Filter:  filter,
	})
	if err != nil {
		return nil, err
	}
	rr := SingleResult(*res)
	return &rr, err
}

func (r *Repository) Find(ctx context.Context, filter *filter.Filter) (*ListResult, error) {
	res, err := r.ds.Find(datasource.QueryContext{
		Context: ctx,
		Filter:  filter,
	})
	if err != nil {
		return nil, err
	}
	rrdata := make([]map[string]interface{}, 0)
	for _, r := range res.Data {
		rrdata = append(rrdata, map[string]interface{}(r))
	}
	return &ListResult{
		Data:  rrdata,
		Count: res.Count,
	}, nil
}

func (r *Repository) UpdateById(ctx context.Context, id interface{}, data map[string]interface{}) (*SingleResult, error) {
	where := map[string]interface{}{
		"id": map[string]interface{}{
			"eq": id,
		},
	}
	_, err := r.ds.Update(datasource.QueryContext{
		Context: ctx,
		Filter: &filter.Filter{
			Where: &where,
		},
	}, data)
	if err != nil {
		return nil, err
	}
	return r.FindById(ctx, id)
}

func (r *Repository) Update(ctx context.Context, filter *filter.Filter, data map[string]interface{}) (*ManyAffectedResult, error) {
	res, err := r.ds.Update(datasource.QueryContext{
		Context: ctx,
		Filter:  filter,
	}, data)
	if err != nil {
		return nil, err
	}
	return &ManyAffectedResult{
		TotalAffected: res.TotalAffected,
	}, nil
}

func (r *Repository) DestroyById(ctx context.Context, id interface{}) error {
	where := map[string]interface{}{
		"id": map[string]interface{}{
			"eq": id,
		},
	}
	_, err := r.ds.Destroy(datasource.QueryContext{
		Context: ctx,
		Filter: &filter.Filter{
			Where: &where,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Destroy(ctx context.Context, filter *filter.Filter) (*ManyAffectedResult, error) {
	res, err := r.ds.Destroy(datasource.QueryContext{
		Context: ctx,
		Filter:  filter,
	})
	if err != nil {
		return nil, err
	}
	return &ManyAffectedResult{
		TotalAffected: res.TotalAffected,
	}, nil
}
