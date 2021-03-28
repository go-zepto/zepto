package linkeradmin

import "github.com/go-zepto/zepto/plugins/linkeradmin/fields"

type Guesser interface {
	Resources() []*Resource
	Icon(resourceName string) string
	ListFields(resourceName string) []*fields.Field
	CreateInputs(resourceName string) []*fields.Input
	UpdateInputs(resourceName string) []*fields.Input
}
