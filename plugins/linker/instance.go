package linker

import (
	"github.com/go-zepto/zepto/plugins/linker/repository"
	"github.com/go-zepto/zepto/web"
)

type LinkerInstance interface {
	Repository(name string) *repository.Repository
	RepositoryDecoder(name string) *repository.RepositoryDecoder
}

func InstanceFromCtx(ctx web.Context) LinkerInstance {
	i := ctx.PluginInstance("linker")
	linkerInstance, ok := i.(LinkerInstance)
	if !ok {
		return nil
	}
	return linkerInstance
}
