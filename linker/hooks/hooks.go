package hooks

import "github.com/go-zepto/zepto/web"

type RemoteHooksInfo struct {
	Endpoint string
	Ctx      web.Context
	ID       *string
	Data     *map[string]interface{}
}

type RemoteHooks interface {
	BeforeRemote(info RemoteHooksInfo) error
	AfterRemote(info RemoteHooksInfo) (*map[string]interface{}, error)
}

type OperationHooksInfo struct {
	Operation string
	ID        *string
	Data      *map[string]interface{}
}

type OperationHooks interface {
	BeforeOperation(info OperationHooksInfo) error
	AfterOperation(info OperationHooksInfo) error
}
