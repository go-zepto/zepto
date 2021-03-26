package linker

import (
	"github.com/gertd/go-pluralize"
	"github.com/go-zepto/zepto/plugins/linker/datasource"
	"github.com/go-zepto/zepto/plugins/linker/utils"
	"github.com/go-zepto/zepto/web"
)

type Linker struct {
	repositories map[string]*Repository
	router       *web.Router
	pluralize    *pluralize.Client
}

func NewLinker(router *web.Router) *Linker {
	return &Linker{
		router:       router,
		repositories: make(map[string]*Repository),
		pluralize:    pluralize.NewClient(),
	}
}

type Resource struct {
	Name           string
	Datasource     datasource.Datasource
	RemoteHooks    RemoteHooks
	OperationHooks OperationHooks
}

func (l *Linker) AddResource(res Resource) {
	if res.RemoteHooks == nil {
		res.RemoteHooks = &DefaultRemoteHooks{}
	}
	if res.OperationHooks == nil {
		res.OperationHooks = &DefaultOperationHooks{}
	}
	repo := NewRepository(RepositoryConfig{
		ResourceName:   res.Name,
		Linker:         l,
		Datasource:     res.Datasource,
		OperationHooks: res.OperationHooks,
	})
	r := RestResource{
		ResourceName: res.Name,
		Linker:       l,
		Repository:   repo,
		RemoteHooks:  res.RemoteHooks,
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

func (l *Linker) Repository(name string) *Repository {
	return l.repositories[name]
}

func (l *Linker) RepositoryDecoder(name string) *RepositoryDecoder {
	return &RepositoryDecoder{
		Repo: l.repositories[name],
	}
}
