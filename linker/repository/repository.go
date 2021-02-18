package repository

import (
	"context"
	"fmt"

	"github.com/go-zepto/zepto/linker/datasource"
	"github.com/go-zepto/zepto/linker/filter"
	"github.com/go-zepto/zepto/linker/hooks"
	"github.com/go-zepto/zepto/linker/utils"
)

type Repository struct {
	ds             datasource.Datasource
	operationHooks hooks.OperationHooks
}

type RepositoryConfig struct {
	Datasource     datasource.Datasource
	OperationHooks hooks.OperationHooks
}

func NewRepository(config RepositoryConfig) *Repository {
	if config.OperationHooks == nil {
		config.OperationHooks = &hooks.DefaultOperationHooks{}
	}
	return &Repository{
		ds:             config.Datasource,
		operationHooks: config.OperationHooks,
	}
}

func (r *Repository) encodeInputDataToMap(input interface{}) map[string]interface{} {
	return utils.DecodeStructToMap(input)
}

func (r *Repository) FindById(ctx context.Context, id interface{}) (*SingleResult, error) {
	where := map[string]interface{}{
		"id": map[string]interface{}{
			"eq": id,
		},
	}
	return r.FindOne(ctx, &filter.Filter{
		Where: &where,
	})
}

func (r *Repository) FindOne(ctx context.Context, filter *filter.Filter) (*SingleResult, error) {
	qc := datasource.QueryContext{
		Context: ctx,
		Filter:  filter,
	}
	err := r.operationHooks.BeforeOperation(hooks.OperationHooksInfo{
		Operation:    "FindOne",
		QueryContext: &qc,
	})
	if err != nil {
		return nil, err
	}
	res, err := r.ds.FindOne(qc)
	if err != nil {
		return nil, err
	}
	rres := *res
	id := fmt.Sprintf("%v", rres["id"])
	err = r.operationHooks.AfterOperation(hooks.OperationHooksInfo{
		Operation:    "FindOne",
		Data:         res,
		QueryContext: &qc,
		ID:           &id,
	})
	rr := SingleResult(*res)
	return &rr, err
}

func (r *Repository) Find(ctx context.Context, filter *filter.Filter) (*ListResult, error) {
	qc := datasource.QueryContext{
		Context: ctx,
		Filter:  filter,
	}
	err := r.operationHooks.BeforeOperation(hooks.OperationHooksInfo{
		Operation:    "Find",
		QueryContext: &qc,
	})
	if err != nil {
		return nil, err
	}
	res, err := r.ds.Find(qc)
	if err != nil {
		return nil, err
	}
	rrdata := make(ManyResults, 0)
	for _, r := range res.Data {
		sr := SingleResult(r)
		rrdata = append(rrdata, &sr)
	}
	lres := ListResult{
		Data:  rrdata,
		Count: res.Count,
	}
	var hres map[string]interface{}
	lres.Decode(&hres)
	err = r.operationHooks.AfterOperation(hooks.OperationHooksInfo{
		Operation:    "Find",
		Data:         &hres,
		QueryContext: &qc,
	})
	if err != nil {
		return nil, err
	}
	return &ListResult{
		Data:  hres["data"].(ManyResults),
		Count: hres["count"].(int64),
	}, nil
}

func (r *Repository) UpdateById(ctx context.Context, id interface{}, data interface{}) (*SingleResult, error) {
	where := map[string]interface{}{
		"id": map[string]interface{}{
			"eq": id,
		},
	}
	filter := &filter.Filter{
		Where: &where,
	}
	_, err := r.Update(ctx, filter, data)
	if err != nil {
		return nil, err
	}
	return r.FindById(ctx, id)
}

func (r *Repository) Update(ctx context.Context, filter *filter.Filter, data interface{}) (*ManyAffectedResult, error) {
	dataMap := r.encodeInputDataToMap(data)
	qc := datasource.QueryContext{
		Context: ctx,
		Filter:  filter,
	}
	err := r.operationHooks.BeforeOperation(hooks.OperationHooksInfo{
		Operation:    "Update",
		QueryContext: &qc,
		Data:         &dataMap,
	})
	if err != nil {
		return nil, err
	}
	res, err := r.ds.Update(qc, data)
	if err != nil {
		return nil, err
	}
	hres := map[string]interface{}{
		"total_affected": res.TotalAffected,
	}
	err = r.operationHooks.AfterOperation(hooks.OperationHooksInfo{
		Operation:    "Update",
		Data:         &hres,
		QueryContext: &qc,
	})
	if err != nil {
		return nil, err
	}
	return &ManyAffectedResult{
		TotalAffected: hres["total_affected"].(int64),
	}, nil
}

func (r *Repository) Create(ctx context.Context, data interface{}) (*SingleResult, error) {
	dataMap := r.encodeInputDataToMap(data)
	qc := datasource.QueryContext{
		Context: ctx,
	}
	err := r.operationHooks.BeforeOperation(hooks.OperationHooksInfo{
		Operation:    "Create",
		QueryContext: &qc,
		Data:         &dataMap,
	})
	if err != nil {
		return nil, err
	}
	res, err := r.ds.Create(qc, dataMap)
	if err != nil {
		return nil, err
	}
	err = r.operationHooks.AfterOperation(hooks.OperationHooksInfo{
		Operation:    "Create",
		Data:         res,
		QueryContext: &qc,
	})
	if err != nil {
		return nil, err
	}
	rr := SingleResult(*res)
	return &rr, nil
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
	qc := datasource.QueryContext{
		Context: ctx,
		Filter:  filter,
	}
	err := r.operationHooks.BeforeOperation(hooks.OperationHooksInfo{
		Operation:    "Destroy",
		QueryContext: &qc,
	})
	if err != nil {
		return nil, err
	}
	res, err := r.ds.Destroy(qc)
	if err != nil {
		return nil, err
	}
	hres := map[string]interface{}{
		"total_affected": res.TotalAffected,
	}
	err = r.operationHooks.AfterOperation(hooks.OperationHooksInfo{
		Operation:    "Destroy",
		Data:         &hres,
		QueryContext: &qc,
	})
	if err != nil {
		return nil, err
	}
	return &ManyAffectedResult{
		TotalAffected: hres["total_affected"].(int64),
	}, nil
}
