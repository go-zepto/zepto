package storage

import (
	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/web"
)

type Options struct {
	Storage Storage
}

func NewStoragePlugin(opts Options) *StoragePlugin {
	return &StoragePlugin{
		storage: opts.Storage,
	}
}

type StoragePlugin struct {
	storage Storage
}

func (s *StoragePlugin) Name() string {
	return "storage"
}

func (s *StoragePlugin) Instance() interface{} {
	return StorageInstance(s.storage)
}

func (s *StoragePlugin) PrependMiddlewares() []web.MiddlewareFunc {
	return nil
}

func (s *StoragePlugin) AppendMiddlewares() []web.MiddlewareFunc {
	return nil
}

func (s *StoragePlugin) OnCreated(z *zepto.Zepto) {}

func (s *StoragePlugin) OnZeptoInit(z *zepto.Zepto) {}

func (s *StoragePlugin) OnZeptoStart(z *zepto.Zepto) {}

func (s *StoragePlugin) OnZeptoStop(z *zepto.Zepto) {}
