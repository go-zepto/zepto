package linker

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-zepto/zepto/plugins/linker/utils"
	"github.com/go-zepto/zepto/web"
)

var ErrUnsuportedMediaType = errors.New(http.StatusText(http.StatusUnsupportedMediaType))

type RestResource struct {
	ResourceName string
	Linker       LinkerInstance
	Repository   *Repository
	RemoteHooks  RemoteHooks
}

func decodeDataFromBody(ctx web.Context, out *map[string]interface{}) error {
	ct := ctx.Request().Header.Get("Content-Type")
	switch ct {
	case "application/json":
		json.NewDecoder(ctx.Request().Body).Decode(&out)
	case "multipart/form-data":
		return nil
	}
	return ErrUnsuportedMediaType
}

func (rest *RestResource) List(ctx web.Context) error {
	err := rest.RemoteHooks.BeforeRemote(RemoteHooksInfo{
		Linker:       rest.Linker,
		ResourceName: rest.ResourceName,
		ResourceID:   nil,
		Endpoint:     "List",
		Ctx:          ctx,
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
	err = rest.RemoteHooks.AfterRemote(RemoteHooksInfo{
		ResourceName: rest.ResourceName,
		Linker:       rest.Linker,
		ResourceID:   nil,
		Endpoint:     "List",
		Result:       &hres,
		Ctx:          ctx,
	})
	if err != nil {
		return err
	}
	return ctx.RenderJson(hres)
}

func (rest *RestResource) Show(ctx web.Context) error {
	id := ctx.Params()["id"]
	err := rest.RemoteHooks.BeforeRemote(RemoteHooksInfo{
		Linker:       rest.Linker,
		ResourceName: rest.ResourceName,
		ResourceID:   &id,
		Data:         nil,
		Endpoint:     "Show",
		Ctx:          ctx,
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
	err = rest.RemoteHooks.AfterRemote(RemoteHooksInfo{
		Linker:       rest.Linker,
		ResourceName: rest.ResourceName,
		ResourceID:   &id,
		Endpoint:     "Show",
		Result:       &hres,
		Ctx:          ctx,
	})
	if err != nil {
		return err
	}
	return ctx.RenderJson(hres)
}

func (rest *RestResource) Create(ctx web.Context) error {
	data := make(map[string]interface{})
	decodeDataFromBody(ctx, &data)
	err := rest.RemoteHooks.BeforeRemote(RemoteHooksInfo{
		Linker:       rest.Linker,
		ResourceName: rest.ResourceName,
		ResourceID:   nil,
		Data:         &data,
		Endpoint:     "Create",
		Ctx:          ctx,
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
	err = rest.RemoteHooks.AfterRemote(RemoteHooksInfo{
		Linker:       rest.Linker,
		ResourceName: rest.ResourceName,
		ResourceID:   &id,
		Endpoint:     "Create",
		Data:         &data,
		Result:       &hres,
		Ctx:          ctx,
	})
	if err != nil {
		return err
	}
	return ctx.RenderJson(hres)
}

func (rest *RestResource) Update(ctx web.Context) error {
	id := ctx.Params()["id"]
	data := make(map[string]interface{})
	decodeDataFromBody(ctx, &data)
	err := rest.RemoteHooks.BeforeRemote(RemoteHooksInfo{
		Linker:       rest.Linker,
		ResourceName: rest.ResourceName,
		ResourceID:   &id,
		Data:         &data,
		Endpoint:     "Update",
		Ctx:          ctx,
	})
	if err != nil {
		return err
	}
	res, err := rest.Repository.UpdateById(ctx, id, data)
	if err != nil {
		return err
	}
	hres := map[string]interface{}(*res)
	err = rest.RemoteHooks.AfterRemote(RemoteHooksInfo{
		Linker:       rest.Linker,
		ResourceName: rest.ResourceName,
		ResourceID:   &id,
		Endpoint:     "Update",
		Data:         &data,
		Result:       &hres,
		Ctx:          ctx,
	})
	if err != nil {
		return err
	}
	return ctx.RenderJson(res)
}

func (rest *RestResource) Destroy(ctx web.Context) error {
	id := ctx.Params()["id"]
	err := rest.RemoteHooks.BeforeRemote(RemoteHooksInfo{
		Linker:       rest.Linker,
		ResourceName: rest.ResourceName,
		ResourceID:   &id,
		Data:         nil,
		Endpoint:     "Destroy",
		Ctx:          ctx,
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
	err = rest.RemoteHooks.AfterRemote(RemoteHooksInfo{
		ResourceName: rest.ResourceName,
		Linker:       rest.Linker,
		ResourceID:   &id,
		Endpoint:     "Destroy",
		Data:         nil,
		Result:       &res,
		Ctx:          ctx,
	})
	if err != nil {
		return err
	}
	return ctx.RenderJson(res)
}
