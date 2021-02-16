package linker

import (
	"github.com/gertd/go-pluralize"
	"github.com/go-zepto/zepto/linker/datasource"
	"github.com/go-zepto/zepto/linker/hooks"
	"github.com/go-zepto/zepto/linker/repository"
	"github.com/go-zepto/zepto/linker/rest"
	"github.com/go-zepto/zepto/linker/utils"
	"github.com/go-zepto/zepto/web"
)

type Linker struct {
	repositories map[string]*repository.Repository
	router       *web.Router
	pluralize    *pluralize.Client
}

func NewLinker(router *web.Router) *Linker {
	return &Linker{
		router:       router,
		repositories: make(map[string]*repository.Repository),
		pluralize:    pluralize.NewClient(),
	}
}

type Resource struct {
	Name           string
	Datasource     datasource.Datasource
	RemoteHooks    hooks.RemoteHooks
	OperationHooks hooks.OperationHooks
}

func (l *Linker) AddResource(res Resource) {
	if res.RemoteHooks == nil {
		res.RemoteHooks = &hooks.DefaultRemoteHooks{}
	}
	if res.OperationHooks == nil {
		res.OperationHooks = &hooks.DefaultOperationHooks{}
	}
	repo := repository.NewRepository(repository.RepositoryConfig{
		Datasource: res.Datasource,
	})
	r := rest.RestResource{
		Repository:     repo,
		RemoteHooks:    res.RemoteHooks,
		OperationHooks: res.OperationHooks,
	}
	l.repositories[res.Name] = repo
	endpoint := utils.ToSnakeCase(l.pluralize.Plural(res.Name))
	l.router.Resource("/"+endpoint, &r)
}

func (l *Linker) Repository(name string) *repository.Repository {
	return l.repositories[name]
}
