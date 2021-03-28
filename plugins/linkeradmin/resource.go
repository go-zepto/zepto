package linkeradmin

import "github.com/go-zepto/zepto/plugins/linker"

type Resource struct {
	Name     string                 `json:"name"`
	Icon     string                 `json:"icon"`
	Endpoint string                 `json:"endpoint"`
	List     *ResourceFieldEndpoint `json:"list_endpoint"`
	Create   *ResourceInputEndpoint `json:"create_endpoint"`
	Update   *ResourceInputEndpoint `json:"update_endpoint"`
}

func NewResource(name string) *Resource {
	return &Resource{
		Name:     name,
		Icon:     "",
		Endpoint: linker.EndpointPathFromResource(name),
		List:     NewResourceFieldEndpoint(),
		Create:   NewResourceInputEndpoint(),
		Update:   NewResourceInputEndpoint(),
	}
}
