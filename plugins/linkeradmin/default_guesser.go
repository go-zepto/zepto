package linkeradmin

import (
	"fmt"

	"github.com/go-zepto/zepto/plugins/linker"
	"github.com/go-zepto/zepto/plugins/linkeradmin/fields"
	"go.uber.org/thriftrw/ptr"
)

// fieldTitleNames is the list of names to be used as title in Admin sorted by priority
var fieldTitleNames = []string{"title", "name", "first_name", "username", "email"}

type DefaultGuesserOptions struct {
	Linker *linker.Linker
}

type DefaultGuesser struct {
	l         *linker.Linker
	resources []*Resource
}

func NewDefaultGuesser(opts DefaultGuesserOptions) *DefaultGuesser {
	resources := make([]*Resource, 0)
	for name := range opts.Linker.Repositories() {
		res := NewResource(name)
		resources = append(resources, res)
	}
	return &DefaultGuesser{
		l:         opts.Linker,
		resources: resources,
	}
}

func (g *DefaultGuesser) Resources() []*Resource {
	return g.resources
}

func (g *DefaultGuesser) getRepoFieldsMap(repo *linker.Repository) map[string]linker.RepositoryField {
	repoFieldsMap := make(map[string]linker.RepositoryField)
	for _, rf := range repo.Fields() {
		repoFieldsMap[rf.Name] = rf
	}
	return repoFieldsMap
}

func (g *DefaultGuesser) guessTitleFieldName(repoFieldsMap map[string]linker.RepositoryField) (string, error) {
	for _, possibleTitleName := range fieldTitleNames {
		if _, ok := repoFieldsMap[possibleTitleName]; ok {
			return possibleTitleName, nil
		}
	}
	return "", fmt.Errorf("could not find possible title name")
}

func (g *DefaultGuesser) ListFields(resourceName string) []*fields.Field {
	repo := g.l.Repository(resourceName)
	if repo == nil {
		return nil
	}
	listFields := make([]*fields.Field, 0)
	repoFieldsMap := g.getRepoFieldsMap(repo)

	// id
	idField := fields.NewTextField("id", nil)
	listFields = append(listFields, &idField)

	// title (best match)
	if titleFieldName, err := g.guessTitleFieldName(repoFieldsMap); err == nil {
		titleField := fields.NewTextField(titleFieldName, nil)
		listFields = append(listFields, &titleField)
	}

	// timestamps (created_at and updated_at)
	if _, ok := repoFieldsMap["created_at"]; ok {
		createdAtField := fields.NewDatetimeField("created_at", &fields.DatetimeFieldOptions{
			ShowTime: ptr.Bool(true),
		})
		listFields = append(listFields, &createdAtField)
	}
	if _, ok := repoFieldsMap["updated_at"]; ok {
		updatedAtField := fields.NewDatetimeField("updated_at", &fields.DatetimeFieldOptions{
			ShowTime: ptr.Bool(true),
		})
		listFields = append(listFields, &updatedAtField)
	}

	return listFields
}

func (g *DefaultGuesser) CreateInputs(resourceName string) []*fields.Input {
	panic("not implemented")
}

func (g *DefaultGuesser) UpdateInputs(resourceName string) []*fields.Input {
	panic("not implemented")
}
