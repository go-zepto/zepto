package rest

import (
	"encoding/json"
	"fmt"

	"github.com/go-zepto/zepto/linker/hooks"
	"github.com/go-zepto/zepto/linker/repository"
	"github.com/go-zepto/zepto/linker/utils"
	"github.com/go-zepto/zepto/web"
)

type RestResource struct {
	Repository     *repository.Repository
	RemoteHooks    hooks.RemoteHooks
	OperationHooks hooks.OperationHooks
}

func (rest *RestResource) List(ctx web.Context) error {
	err := rest.RemoteHooks.BeforeRemote(hooks.RemoteHooksInfo{
		ID:       nil,
		Endpoint: "List",
		Ctx:      ctx,
	})
	if err != nil {
		return err
	}
	res, err := rest.Repository.Find(ctx, utils.GetFilterFromQueryArgCtx(ctx))
	if err != nil {
		return err
	}
	var hres map[string]interface{}
	res.DecodeAll(&hres)
	ares, err := rest.RemoteHooks.AfterRemote(hooks.RemoteHooksInfo{
		ID:       nil,
		Endpoint: "List",
		Data:     &hres,
		Ctx:      ctx,
	})
	if err != nil {
		return err
	}
	return ctx.RenderJson(ares)
}

func (rest *RestResource) Show(ctx web.Context) error {
	id := ctx.Params()["id"]
	err := rest.RemoteHooks.BeforeRemote(hooks.RemoteHooksInfo{
		ID:       &id,
		Data:     nil,
		Endpoint: "Show",
		Ctx:      ctx,
	})
	if err != nil {
		return err
	}
	res, err := rest.Repository.FindById(ctx, id)
	if err != nil {
		ctx.SetStatus(400)
		return err
	}
	var hres map[string]interface{}
	res.Decode(&hres)
	ares, err := rest.RemoteHooks.AfterRemote(hooks.RemoteHooksInfo{
		ID:       &id,
		Endpoint: "Show",
		Data:     &hres,
		Ctx:      ctx,
	})
	if err != nil {
		return err
	}
	return ctx.RenderJson(ares)
}

func (rest *RestResource) Create(ctx web.Context) error {
	var data map[string]interface{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&data)
	if err != nil {
		return err
	}
	err = rest.RemoteHooks.BeforeRemote(hooks.RemoteHooksInfo{
		ID:       nil,
		Data:     &data,
		Endpoint: "Create",
		Ctx:      ctx,
	})
	if err != nil {
		return err
	}
	res, err := rest.Repository.Create(ctx, data)
	if err != nil {
		return err
	}
	var hres map[string]interface{}
	res.Decode(&hres)
	id := fmt.Sprintf("%v", hres["id"])
	ares, err := rest.RemoteHooks.AfterRemote(hooks.RemoteHooksInfo{
		ID:       &id,
		Endpoint: "Create",
		Data:     &hres,
		Ctx:      ctx,
	})
	if err != nil {
		return err
	}
	return ctx.RenderJson(ares)
}

func (rest *RestResource) Update(ctx web.Context) error {
	id := ctx.Params()["id"]
	var data map[string]interface{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&data)
	if err != nil {
		return err
	}
	err = rest.RemoteHooks.BeforeRemote(hooks.RemoteHooksInfo{
		ID:       &id,
		Data:     &data,
		Endpoint: "Update",
		Ctx:      ctx,
	})
	if err != nil {
		return err
	}
	res, err := rest.Repository.UpdateById(ctx, id, data)
	if err != nil {
		return err
	}
	var hres map[string]interface{}
	res.Decode(&hres)
	ares, err := rest.RemoteHooks.AfterRemote(hooks.RemoteHooksInfo{
		ID:       &id,
		Endpoint: "Update",
		Data:     &hres,
		Ctx:      ctx,
	})
	if err != nil {
		return err
	}
	return ctx.RenderJson(ares)
}

func (rest *RestResource) Destroy(ctx web.Context) error {
	id := ctx.Params()["id"]
	err := rest.RemoteHooks.BeforeRemote(hooks.RemoteHooksInfo{
		ID:       &id,
		Data:     nil,
		Endpoint: "Destroy",
		Ctx:      ctx,
	})
	if err != nil {
		return err
	}
	err = rest.Repository.DestroyById(ctx, id)
	if err != nil {
		return err
	}
	res := map[string]interface{}{
		"deleted": true,
	}
	ares, err := rest.RemoteHooks.AfterRemote(hooks.RemoteHooksInfo{
		ID:       &id,
		Endpoint: "Destroy",
		Data:     &res,
		Ctx:      ctx,
	})
	if err != nil {
		return err
	}
	return ctx.RenderJson(ares)
}
