package linkeradmin

import "github.com/go-zepto/zepto/plugins/linker"

type LinkerResource struct {
	Name     string                 `json:"name"`
	Endpoint string                 `json:"endpoint"`
	List     *ResourceFieldEndpoint `json:"list_endpoint"`
	Create   *ResourceFieldEndpoint `json:"create_endpoint"`
	Update   *ResourceFieldEndpoint `json:"update_endpoint"`
}

func NewLinkerResource(name string) *LinkerResource {
	return &LinkerResource{
		Name:     name,
		Endpoint: linker.EndpointPathFromResource(name),
		List:     NewResourceFieldEndpoint(),
	}
}
