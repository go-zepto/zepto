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

func TestGuessCreateInputs(t *testing.T) {
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

	createInputs := guesser.CreateInputs("Book")

	expectedInputs := []fields.Input{
		fields.NewTextInput("name", nil),
	}
	assert.Len(t, expectedInputs, len(createInputs))
	for idx, ef := range expectedInputs {
		assert.Equal(t, &ef, createInputs[idx])
	}
}

func TestGuessCreateInputs_Datetime(t *testing.T) {
	t.Skip("Unstable test")
	guesser := newTestDefaultGuesser(t, map[string]datasource.Field{
		"id": {
			Name: "id",
			Type: "integer",
		},
		"name": {
			Name: "name",
			Type: "text",
		},
		"expiration": {
			Name: "expiration",
			Type: "datetime",
		},
		"active_at": {
			Name: "active_at",
			Type: "datetime",
		},
	})

	createInputs := guesser.CreateInputs("Book")

	expectedInputs := []fields.Input{
		fields.NewTextInput("name", nil),
		fields.NewDatetimeInput("expiration", nil),
		fields.NewDatetimeInput("active_at", nil),
	}
	assert.Len(t, expectedInputs, len(createInputs))
	for idx, ef := range expectedInputs {
		assert.Equal(t, &ef, createInputs[idx])
	}
}

func TestGuessCreateInputs_ReferenceInput(t *testing.T) {

	ctrl := gomock.NewController(t)
	ds := datasourcemock.NewMockDatasource(ctrl)
	ds.EXPECT().
		Fields().Return(map[string]datasource.Field{})

	guesser := newTestDefaultGuesser(t, map[string]datasource.Field{
		"id": {
			Name: "id",
			Type: "integer",
		},
		"name": {
			Name: "name",
			Type: "text",
		},
		"author_id": {
			Name: "author_id",
			Type: "numeric",
		},
	})

	guesser.l.AddResource(linker.Resource{
		Name:       "Author",
		Datasource: ds,
	})

	createInputs := guesser.CreateInputs("Book")

	expectedInputs := []fields.Input{
		fields.NewTextInput("name", nil),
		fields.NewReferenceInput("author_id", "Author", &fields.ReferenceInputOptions{
			Autocomplete: fields.ReferenceInputAutocomplete{
				Enabled: true,
			},
		}),
	}
	assert.Len(t, expectedInputs, len(createInputs))
	for idx, ef := range expectedInputs {
		assert.Equal(t, &ef, createInputs[idx])
	}
}
