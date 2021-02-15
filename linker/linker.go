package linker

import (
	"github.com/gertd/go-pluralize"
	"github.com/go-zepto/zepto/linker/datasource"
	"github.com/go-zepto/zepto/linker/repository"
	"github.com/go-zepto/zepto/linker/rest"
	"github.com/go-zepto/zepto/linker/utils"
	"github.com/go-zepto/zepto/web"
)

type Linker struct {
	resources map[string]rest.RestResource
	router    *web.Router
	pluralize *pluralize.Client
}

func NewLinker(router *web.Router) *Linker {
	return &Linker{
		router:    router,
		resources: make(map[string]rest.RestResource),
		pluralize: pluralize.NewClient(),
	}
}

func (l *Linker) AddResource(name string, ds datasource.Datasource) {
	r := rest.RestResource{
		Repository: repository.NewRepository(repository.RepositoryConfig{
			Datasource: ds,
		}),
	}
	endpoint := utils.ToSnakeCase(l.pluralize.Plural(name))
	l.router.Resource("/"+endpoint, &r)
}
