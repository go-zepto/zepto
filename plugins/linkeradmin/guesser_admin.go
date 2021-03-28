package linkeradmin

import (
	"github.com/go-zepto/zepto/plugins/linker"
)

// fieldTitleNames is the prioritized list of names to be used as title in Admin
var fieldTitleNames = []string{"title", "name", "first_name", "username", "email"}

var fieldTypeAliases = map[string][]string{
	"number": {
		"integer",
		"float",
	},
}

type linkerFieldMap map[string]linker.RepositoryField

func NewGuesserAdmin(l *linker.Linker) *Admin {
	admin := NewAdmin()
	for name, repo := range l.Repositories() {
		res := NewLinkerResource(name)
		fieldMap := make(linkerFieldMap)
		for _, f := range repo.Fields() {
			fieldMap[f.Name] = f
		}
		guessList(res, fieldMap)
		guessCreate(res, fieldMap)
		admin.AddResource(res)
	}
	return admin
}

func findFirstFieldFromPriotizedList(fieldMap linkerFieldMap, list []string) *linker.RepositoryField {
	for _, fieldName := range list {
		f, ok := fieldMap[fieldName]
		if ok {
			return &f
		}
	}
	return nil
}

func findTimestamps(fieldMap linkerFieldMap) []*linker.RepositoryField {
	timestampsFields := make([]*linker.RepositoryField, 0)
	for _, fieldName := range [...]string{"created_at", "updated_at"} {
		tsField, ok := fieldMap[fieldName]
		if ok {
			timestampsFields = append(timestampsFields, &tsField)
		}
	}
	return timestampsFields
}

func guessList(res *LinkerResource, fieldMap linkerFieldMap) {
	id, ok := fieldMap["id"]
	if ok {
		res.List.AddField(Field{
			Name: id.Name,
			Type: "text",
		})
	}
	title := findFirstFieldFromPriotizedList(fieldMap, fieldTitleNames)
	if title != nil {
		res.List.AddField(Field{
			Name: title.Name,
			Type: "text",
		})
	}
	timestamps := findTimestamps(fieldMap)
	for _, ts := range timestamps {
		res.List.AddField(Field{
			Name: ts.Name,
			Type: "datetime",
			Options: map[string]interface{}{
				"props": map[string]interface{}{
					"showTime": true,
				},
			},
		})
	}
}

func guessCreate(res *LinkerResource, fieldMap linkerFieldMap) {
	id, ok := fieldMap["id"]
	if ok {
		res.Create.AddInput(Field{
			Name: id.Name,
			Type: "text",
			Options: map[string]interface{}{
				"props": map[string]interface{}{
					"disabled": true,
				},
			},
		})
	}
	res.Create.AddInput(Field{
		Name: "name",
		Type: "text",
	})
}
