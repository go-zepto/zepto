package linkeradmin

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-zepto/zepto/plugins/linker"
	"github.com/go-zepto/zepto/plugins/linker/utils"
	"github.com/go-zepto/zepto/plugins/linkeradmin/fields"
	"go.uber.org/thriftrw/ptr"
)

// fieldTitleNames is the list of names to be used as title in Admin sorted by priority
var fieldTitleNames = []string{"title", "name", "first_name", "username", "email"}
var commonBooleanNamesRegex = `is_(.*)|(active|enabled|visible|hidden|checked|approved|paused|ready)`

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

func (g *DefaultGuesser) guessReferenceInputOptionText(resourceName string) string {
	refOptionTextField, _ := g.guessTitleFieldName(g.getRepoFieldsMap(g.l.Repository(resourceName)))
	return refOptionTextField
}

func (g *DefaultGuesser) guessReferenceInput(fieldName string) *fields.Input {
	repoNames := make(map[string]string)
	for repoName := range g.l.Repositories() {
		repoNames[strings.ToLower(utils.ToSnakeCase(repoName))] = repoName
	}
	if resourceName, ok := repoNames[strings.TrimSuffix(fieldName, "_id")]; ok {
		refOptionTextField := g.guessReferenceInputOptionText(resourceName)
		refInput := fields.NewReferenceInput(fieldName, resourceName, &fields.ReferenceInputOptions{
			OptionTextField: refOptionTextField,
			Autocomplete: fields.ReferenceInputAutocomplete{
				Enabled: true,
			},
		})
		return &refInput
	}
	return nil
}

func (g *DefaultGuesser) createGuessedInput(rf *linker.RepositoryField) fields.Input {
	// Direct Type checks
	switch rf.Type {
	case "datetime":
		return fields.NewDatetimeInput(rf.Name, nil)
	case "text":
		return fields.NewTextInput(rf.Name, nil)
	}

	// Regex checks
	rfNameBytes := []byte(rf.Name)

	// Boolean
	if isBool, _ := regexp.Match(commonBooleanNamesRegex, rfNameBytes); isBool {
		return fields.NewTextInput(rf.Name, nil)
	}

	// Reference
	if isRef, _ := regexp.Match(`.*_id`, rfNameBytes); isRef {
		refInput := g.guessReferenceInput(rf.Name)
		if refInput != nil {
			return *refInput
		}
	}

	// Falback to Text Input
	return fields.NewTextInput(rf.Name, nil)
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
	repo := g.l.Repository(resourceName)
	if repo == nil {
		return nil
	}
	createInputs := make([]*fields.Input, 0)
	repoFieldsMap := g.getRepoFieldsMap(repo)

	// These fields are not relevant in create form
	delete(repoFieldsMap, "id")
	delete(repoFieldsMap, "created_at")
	delete(repoFieldsMap, "updated_at")
	delete(repoFieldsMap, "deleted_at")

	for _, rf := range repoFieldsMap {
		input := g.createGuessedInput(&rf)
		createInputs = append(createInputs, &input)
	}
	return createInputs
}

func (g *DefaultGuesser) UpdateInputs(resourceName string) []*fields.Input {
	repo := g.l.Repository(resourceName)
	if repo == nil {
		return nil
	}
	createInputs := make([]*fields.Input, 0)
	repoFieldsMap := g.getRepoFieldsMap(repo)

	// Removing to create disabled field
	delete(repoFieldsMap, "id")
	delete(repoFieldsMap, "created_at")
	delete(repoFieldsMap, "updated_at")
	delete(repoFieldsMap, "deleted_at")

	idInput := fields.NewTextInput("id", &fields.TextInputOptions{
		Disabled: ptr.Bool(true),
	})
	createInputs = append(createInputs, &idInput)

	for _, rf := range repoFieldsMap {
		input := g.createGuessedInput(&rf)
		createInputs = append(createInputs, &input)
	}

	// timestamps (created_at and updated_at)
	if _, ok := repoFieldsMap["created_at"]; ok {
		createdAtInput := fields.NewDatetimeInput("created_at", &fields.DatetimeInputOptions{
			Disabled: ptr.Bool(true),
		})
		createInputs = append(createInputs, &createdAtInput)
	}
	if _, ok := repoFieldsMap["updated_at"]; ok {
		updatedAtInput := fields.NewDatetimeInput("updated_at", &fields.DatetimeInputOptions{
			Disabled: ptr.Bool(true),
		})
		createInputs = append(createInputs, &updatedAtInput)
	}

	return createInputs
}
