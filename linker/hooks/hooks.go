package hooks

import (
	"github.com/go-zepto/zepto/linker/datasource"
	"github.com/go-zepto/zepto/web"
)

type RemoteHooksInfo struct {
	Endpoint string
	Ctx      web.Context
	ID       *string
	Data     *map[string]interface{}
}

type RemoteHooks interface {
	BeforeRemote(info RemoteHooksInfo) error
	AfterRemote(info RemoteHooksInfo) error
}

type OperationHooksInfo struct {
	Operation    string
	ID           *string
	Data         *map[string]interface{}
	QueryContext *datasource.QueryContext
}

type OperationHooks interface {
	BeforeOperation(info OperationHooksInfo) error
	AfterOperation(info OperationHooksInfo) error
}
