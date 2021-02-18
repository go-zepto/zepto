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
	Repository  *repository.Repository
	RemoteHooks hooks.RemoteHooks
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
	filter, err := utils.GetFilterFromQueryArgCtx(ctx)
	if err != nil {
		ctx.SetStatus(400)
		return ctx.RenderJson(map[string]string{
			"error": err.Error(),
		})
	}
	res, err := rest.Repository.Find(ctx, filter)
	if err != nil {
		return err
	}
	var hres map[string]interface{}
	res.Decode(&hres)
	err = rest.RemoteHooks.AfterRemote(hooks.RemoteHooksInfo{
		ID:       nil,
		Endpoint: "List",
		Data:     &hres,
		Ctx:      ctx,
	})
	if err != nil {
		return err
	}
	return ctx.RenderJson(hres)
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
	hres := map[string]interface{}(*res)
	err = rest.RemoteHooks.AfterRemote(hooks.RemoteHooksInfo{
		ID:       &id,
		Endpoint: "Show",
		Data:     &hres,
		Ctx:      ctx,
	})
	if err != nil {
		return err
	}
	return ctx.RenderJson(hres)
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
	hres := map[string]interface{}(*res)
	id := fmt.Sprintf("%v", hres["id"])
	err = rest.RemoteHooks.AfterRemote(hooks.RemoteHooksInfo{
		ID:       &id,
		Endpoint: "Create",
		Data:     &hres,
		Ctx:      ctx,
	})
	if err != nil {
		return err
	}
	return ctx.RenderJson(hres)
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
	hres := map[string]interface{}(*res)
	err = rest.RemoteHooks.AfterRemote(hooks.RemoteHooksInfo{
		ID:       &id,
		Endpoint: "Update",
		Data:     &hres,
		Ctx:      ctx,
	})
	if err != nil {
		return err
	}
	return ctx.RenderJson(res)
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
	err = rest.RemoteHooks.AfterRemote(hooks.RemoteHooksInfo{
		ID:       &id,
		Endpoint: "Destroy",
		Data:     &res,
		Ctx:      ctx,
	})
	if err != nil {
		return err
	}
	return ctx.RenderJson(res)
}
