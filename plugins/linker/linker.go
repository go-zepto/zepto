package linker

import (
	"github.com/gertd/go-pluralize"
	"github.com/go-zepto/zepto/plugins/linker/datasource"
	"github.com/go-zepto/zepto/plugins/linker/hooks"
	"github.com/go-zepto/zepto/plugins/linker/repository"
	"github.com/go-zepto/zepto/plugins/linker/rest"
	"github.com/go-zepto/zepto/plugins/linker/utils"
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
		Datasource:     res.Datasource,
		OperationHooks: res.OperationHooks,
	})
	r := rest.RestResource{
		Repository:  repo,
		RemoteHooks: res.RemoteHooks,
	}
	l.repositories[res.Name] = repo
	endpoint := utils.ToSnakeCase(l.pluralize.Plural(res.Name))
	l.router.Resource("/"+endpoint, &r)
}

func (l *Linker) AddResources(resources []Resource) {
	for _, r := range resources {
		l.AddResource(r)
	}
}

func (l *Linker) Router() *web.Router {
	return l.router
}

func (l *Linker) Repository(name string) *repository.Repository {
	return l.repositories[name]
}

func (l *Linker) RepositoryDecoder(name string) *repository.RepositoryDecoder {
	return &repository.RepositoryDecoder{
		Repo: l.repositories[name],
	}
}
