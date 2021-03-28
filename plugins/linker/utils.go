package linker

import (
	"github.com/gertd/go-pluralize"
	"github.com/go-zepto/zepto/plugins/linker/utils"
)

var p = pluralize.NewClient()

func EndpointPathFromResource(resourceName string) string {
	return utils.ToSnakeCase(p.Plural(resourceName))
}
