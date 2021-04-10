package linker

import (
	"context"

	"github.com/go-zepto/zepto/plugins/linker/filter"
)

type RepositoryDecoder struct {
	Repo *Repository
}

func (r *RepositoryDecoder) FindById(ctx context.Context, id interface{}, dest interface{}) error {
	res, err := r.Repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	return res.Decode(dest)
}

func (r *RepositoryDecoder) FindOne(ctx context.Context, filter *filter.Filter, dest interface{}) error {
	res, err := r.Repo.FindOne(ctx, filter)
	if err != nil {
		return err
	}
	return res.Decode(dest)
}

func (r *RepositoryDecoder) Find(ctx context.Context, filter *filter.Filter, dest interface{}) error {
	res, err := r.Repo.Find(ctx, filter)
	if err != nil {
		return err
	}
	return res.Decode(dest)
}

func (r *RepositoryDecoder) UpdateById(ctx context.Context, id interface{}, data interface{}, dest interface{}) error {
	res, err := r.Repo.UpdateById(ctx, id, data)
	if err != nil {
		return err
	}
	return res.Decode(dest)
}

func (r *RepositoryDecoder) Update(ctx context.Context, filter *filter.Filter, data interface{}, dest interface{}) error {
	res, err := r.Repo.Update(ctx, filter, data)
	if err != nil {
		return err
	}
	return res.Decode(dest)
}

func (r *RepositoryDecoder) Create(ctx context.Context, data interface{}, dest interface{}) error {
	res, err := r.Repo.Create(ctx, data)
	if err != nil {
		return err
	}
	return res.Decode(dest)
}

func (r *RepositoryDecoder) DestroyById(ctx context.Context, id interface{}) error {
	err := r.Repo.DestroyById(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepositoryDecoder) Destroy(ctx context.Context, filter *filter.Filter, dest interface{}) error {
	res, err := r.Repo.Destroy(ctx, filter)
	if err != nil {
		return err
	}
	return res.Decode(dest)
}
