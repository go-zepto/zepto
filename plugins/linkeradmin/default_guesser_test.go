package linkeradmin

import (
	"testing"

	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/plugins/linker"
	"github.com/go-zepto/zepto/plugins/linker/datasource"
	"github.com/go-zepto/zepto/plugins/linker/datasource/datasourcemock"
	"github.com/go-zepto/zepto/plugins/linkeradmin/fields"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/thriftrw/ptr"
)

func newTestDefaultGuesser(t *testing.T, fields map[string]datasource.Field) *DefaultGuesser {
	ctrl := gomock.NewController(t)
	ds := datasourcemock.NewMockDatasource(ctrl)

	ds.EXPECT().
		Fields().Return(fields)

	z := zepto.NewZepto()
	l := linker.NewLinker(z.Router("/api"))
	l.AddResource(linker.Resource{
		Name:       "Book",
		Datasource: ds,
	})
	return NewDefaultGuesser(DefaultGuesserOptions{
		Linker: l,
	})
}

func TestGuessTitleFieldName(t *testing.T) {
	guesser := newTestDefaultGuesser(t, map[string]datasource.Field{
		"id": {
			Name: "id",
			Type: "integer",
		},
		"name": {
			Name: "name",
			Type: "text",
		},
	})
	repo := guesser.l.Repository("Book")
	fieldName, err := guesser.guessTitleFieldName(guesser.getRepoFieldsMap(repo))
	assert.NoError(t, err)
	assert.Equal(t, "name", fieldName)
}

func TestGuessTitleFieldName_BestMatch(t *testing.T) {
	guesser := newTestDefaultGuesser(t, map[string]datasource.Field{
		"id": {
			Name: "id",
			Type: "integer",
		},
		"name": {
			Name: "name",
			Type: "text",
		},
		"title": {
			Name: "title",
			Type: "text",
		},
		"first_name": {
			Name: "first_name",
			Type: "text",
		},
	})
	repo := guesser.l.Repository("Book")
	fieldName, err := guesser.guessTitleFieldName(guesser.getRepoFieldsMap(repo))
	assert.NoError(t, err)
	assert.Equal(t, "title", fieldName)
}

func TestGuessListFields(t *testing.T) {
	guesser := newTestDefaultGuesser(t, map[string]datasource.Field{
		"id": {
			Name: "id",
			Type: "integer",
		},
		"name": {
			Name: "name",
			Type: "text",
		},
	})

	listFields := guesser.ListFields("Book")

	expectedFields := []fields.Field{
		fields.NewTextField("id", nil),
		fields.NewTextField("name", nil),
	}
	assert.Len(t, expectedFields, len(listFields))
	for idx, ef := range expectedFields {
		assert.Equal(t, &ef, listFields[idx])
	}
}

func TestGuessListFields_Timestamps(t *testing.T) {
	guesser := newTestDefaultGuesser(t, map[string]datasource.Field{
		"id": {
			Name: "id",
			Type: "integer",
		},
		"name": {
			Name: "name",
			Type: "text",
		},
		"created_at": {
			Name: "created_at",
			Type: "datetime",
		},
		"updated_at": {
			Name: "updated_at",
			Type: "datetime",
		},
	})

	listFields := guesser.ListFields("Book")

	expectedFields := []fields.Field{
		fields.NewTextField("id", nil),
		fields.NewTextField("name", nil),
		fields.NewDatetimeField("created_at", &fields.DatetimeFieldOptions{
			ShowTime: ptr.Bool(true),
		}),
		fields.NewDatetimeField("updated_at", &fields.DatetimeFieldOptions{
			ShowTime: ptr.Bool(true),
		}),
	}
	assert.Len(t, expectedFields, len(listFields))
	for idx, ef := range expectedFields {
		assert.Equal(t, &ef, listFields[idx])
	}
}
