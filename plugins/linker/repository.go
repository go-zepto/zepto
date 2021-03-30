package linker

import (
	"context"
	"fmt"

	"github.com/go-zepto/zepto/plugins/linker/datasource"
	"github.com/go-zepto/zepto/plugins/linker/filter"
	"github.com/go-zepto/zepto/plugins/linker/utils"
)

type Repository struct {
	resourceName   string
	linker         LinkerInstance
	ds             datasource.Datasource
	operationHooks OperationHooks
}

type RepositoryConfig struct {
	ResourceName   string
	Linker         LinkerInstance
	Datasource     datasource.Datasource
	OperationHooks OperationHooks
}

type RepositoryField struct {
	Name     string
	Type     string
	Nullable bool
}

func NewRepository(config RepositoryConfig) *Repository {
	if config.OperationHooks == nil {
		config.OperationHooks = &DefaultOperationHooks{}
	}
	return &Repository{
		resourceName:   config.ResourceName,
		linker:         config.Linker,
		ds:             config.Datasource,
		operationHooks: config.OperationHooks,
	}
}

func ctxWithResourceId(ctx context.Context, id interface{}) context.Context {
	return context.WithValue(ctx, "resource_id", fmt.Sprintf("%v", id))
}

func resourceIdFromCtx(ctx context.Context) *string {
	var resourceId *string
	resId, _ := ctx.Value("resource_id").(string)
	if resId != "" {
		resourceId = &resId
	}
	return resourceId
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
	return r.FindOne(ctxWithResourceId(ctx, id), &filter.Filter{
		Where: &where,
	})
}

func (r *Repository) FindOne(ctx context.Context, filter *filter.Filter) (*SingleResult, error) {
	qc := datasource.QueryContext{
		Context: ctx,
		Filter:  filter,
	}
	err := r.operationHooks.BeforeOperation(OperationHooksInfo{
		ResourceName: r.resourceName,
		ResourceID:   resourceIdFromCtx(ctx),
		Linker:       r.linker,
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
	err = r.operationHooks.AfterOperation(OperationHooksInfo{
		ResourceName: r.resourceName,
		Linker:       r.linker,
		Operation:    "FindOne",
		Result:       res,
		QueryContext: &qc,
		ResourceID:   &id,
	})
	rr := SingleResult(*res)
	return &rr, err
}

func (r *Repository) Find(ctx context.Context, filter *filter.Filter) (*ListResult, error) {
	qc := datasource.QueryContext{
		Context: ctx,
		Filter:  filter,
	}
	err := r.operationHooks.BeforeOperation(OperationHooksInfo{
		ResourceName: r.resourceName,
		Linker:       r.linker,
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
	err = r.operationHooks.AfterOperation(OperationHooksInfo{
		ResourceName: r.resourceName,
		Linker:       r.linker,
		Operation:    "Find",
		Result:       &hres,
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
	_, err := r.Update(ctxWithResourceId(ctx, id), filter, data)
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
	err := r.operationHooks.BeforeOperation(OperationHooksInfo{
		ResourceName: r.resourceName,
		ResourceID:   resourceIdFromCtx(ctx),
		Linker:       r.linker,
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
	err = r.operationHooks.AfterOperation(OperationHooksInfo{
		ResourceName: r.resourceName,
		ResourceID:   resourceIdFromCtx(ctx),
		Linker:       r.linker,
		Operation:    "Update",
		Data:         &dataMap,
		Result:       &hres,
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
	err := r.operationHooks.BeforeOperation(OperationHooksInfo{
		ResourceName: r.resourceName,
		Linker:       r.linker,
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
	rres := *res
	id := fmt.Sprintf("%v", rres["id"])
	err = r.operationHooks.AfterOperation(OperationHooksInfo{
		ResourceName: r.resourceName,
		ResourceID:   &id,
		Linker:       r.linker,
		Operation:    "Create",
		Data:         &dataMap,
		Result:       res,
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
	_, err := r.Destroy(ctxWithResourceId(ctx, id), &filter.Filter{
		Where: &where,
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
	err := r.operationHooks.BeforeOperation(OperationHooksInfo{
		ResourceName: r.resourceName,
		ResourceID:   resourceIdFromCtx(ctx),
		Linker:       r.linker,
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
	err = r.operationHooks.AfterOperation(OperationHooksInfo{
		ResourceName: r.resourceName,
		ResourceID:   resourceIdFromCtx(ctx),
		Linker:       r.linker,
		Operation:    "Destroy",
		Result:       &hres,
		QueryContext: &qc,
	})
	if err != nil {
		return nil, err
	}
	return &ManyAffectedResult{
		TotalAffected: hres["total_affected"].(int64),
	}, nil
}

func (r *Repository) Fields() []RepositoryField {
	fields := make([]RepositoryField, 0)
	for _, f := range r.ds.Fields() {
		fields = append(fields, RepositoryField{
			Name:     f.Name,
			Type:     f.Type,
			Nullable: f.Nullable,
		})
	}
	return fields
}
