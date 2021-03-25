package storage

import "github.com/go-zepto/zepto/web"

type StorageInstance interface {
	Storage
}

func InstanceFromCtx(ctx web.Context) StorageInstance {
	i := ctx.PluginInstance("storage")
	storageInstance, ok := i.(StorageInstance)
	if !ok {
		return nil
	}
	return storageInstance
}
