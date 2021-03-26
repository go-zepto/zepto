package linker

import (
	"github.com/go-zepto/zepto/plugins/linker/datasource"
	"github.com/go-zepto/zepto/web"
)

type RemoteHooksInfo struct {
	// Type of endpoint: List, Show, Create, Update or Destroy
	Endpoint string

	// Zepto Web Context
	Ctx web.Context

	// Current Resource Name
	ResourceName string

	// Current Resource ID (Only present in Show, Update and Destroy endpoints)
	ResourceID *string

	// Data parsed from JSON body (Only useful in Create and Update endpoints)
	Data *map[string]interface{}

	// Result data that will be encoded to JSON and be sent to the user
	Result *map[string]interface{}

	// Linker Plugin Instance. You can use it to do another queries inside hooks
	Linker LinkerInstance
}

/*
	RemoteHooks is used to intercept before and after a HTTP Remote request

	It has access to HTTP request and you can use all Zepto web features like:
	cookies, sessions, headers, plugins, etc.


	Use RemoteHooks to restrict some users to access specifics endpoints or to customize the final body json sent to user.

	When BeforeRemote or AfterRemote returns an error, the linker will not continue and return http error.

	Note: This is not a controller. Although it is possible, we do not recommend calling the render here.
	You can only prevent the request from happening, but if you need to return a body completely different from Linker, it is recommended to create a controller and route using the common Zepto pathways.
*/
type RemoteHooks interface {
	BeforeRemote(info RemoteHooksInfo) error
	AfterRemote(info RemoteHooksInfo) error
}

type OperationHooksInfo struct {
	// Type of operation: FindOne, Find, Update, Create or Destroy
	Operation string

	// Current Resource Name
	ResourceName string

	// Current Resource ID (Only present in FindOne, Update or Destroy operations)
	ResourceID *string

	// Data used to perform the OperationHooks Create and Update
	Data *map[string]interface{}

	// Result from the AfterOperation it can be a SingleResult or a ListResult depending on the operation
	Result *map[string]interface{}

	/*
		QueryContext is the object used to apply the Linker filters.

		Use this object to intercept or change where, limit, skip and include filters before
		perform the OperationHooks.
	*/
	QueryContext *datasource.QueryContext

	// Linker Plugin Instance. You can use it to do another queries inside hooks
	Linker LinkerInstance
}

/*
	OperationHooks is used to intercept before and after an operation to datasource (Gorm, etc)

	It doesn't know what is HTTP and doesn't have context and access to HTTP information.

	Use OperationHooks when you need to intercept a query in all cases, from the default HTTP request
	and from the programmatically call like:

	Linker.Respository("User").Find(...)
*/
type OperationHooks interface {
	BeforeOperation(info OperationHooksInfo) error
	AfterOperation(info OperationHooksInfo) error
}

// Defaults

type DefaultRemoteHooks struct{}

func (drh *DefaultRemoteHooks) BeforeRemote(info RemoteHooksInfo) error {
	return nil
}

func (drh *DefaultRemoteHooks) AfterRemote(info RemoteHooksInfo) error {
	return nil
}

type DefaultOperationHooks struct{}

func (drh *DefaultOperationHooks) BeforeOperation(info OperationHooksInfo) error {
	return nil
}

func (drh *DefaultOperationHooks) AfterOperation(info OperationHooksInfo) error {
	return nil
}
