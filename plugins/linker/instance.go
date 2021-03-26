package linker

import (
	"github.com/go-zepto/zepto/web"
)

type LinkerInstance interface {
	Repository(name string) *Repository
	RepositoryDecoder(name string) *RepositoryDecoder
}

func InstanceFromCtx(ctx web.Context) LinkerInstance {
	i := ctx.PluginInstance("linker")
	linkerInstance, ok := i.(LinkerInstance)
	if !ok {
		return nil
	}
	return linkerInstance
}
