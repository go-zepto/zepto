package rest

import (
	"encoding/json"

	"github.com/go-zepto/zepto/linker/repository"
	"github.com/go-zepto/zepto/linker/utils"
	"github.com/go-zepto/zepto/web"
)

type RestResource struct {
	Repository *repository.Repository
}

func (rest *RestResource) List(ctx web.Context) error {
	res, err := rest.Repository.Find(ctx, utils.GetFilterFromQueryArgCtx(ctx))
	if err != nil {
		return err
	}
	return ctx.RenderJson(res)
}

func (rest *RestResource) Show(ctx web.Context) error {
	id := ctx.Params()["id"]
	res, err := rest.Repository.FindById(ctx, id)
	if err != nil {
		ctx.SetStatus(400)
		return err
	}
	return ctx.RenderJson(res)
}

func (rest *RestResource) Create(ctx web.Context) error {
	var data map[string]interface{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&data)
	if err != nil {
		return err
	}
	res, err := rest.Repository.Create(ctx, data)
	if err != nil {
		return err
	}
	return ctx.RenderJson(res)
}

func (rest *RestResource) Update(ctx web.Context) error {
	id := ctx.Params()["id"]
	var data map[string]interface{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&data)
	if err != nil {
		return err
	}
	res, err := rest.Repository.UpdateById(ctx, id, data)
	if err != nil {
		return err
	}
	return ctx.RenderJson(res)
}

func (rest *RestResource) Destroy(ctx web.Context) error {
	id := ctx.Params()["id"]
	err := rest.Repository.DestroyById(ctx, id)
	if err != nil {
		return err
	}
	return ctx.RenderJson(map[string]bool{
		"deleted": true,
	})
}
