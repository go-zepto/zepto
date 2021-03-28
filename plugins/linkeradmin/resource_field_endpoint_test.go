package linkeradmin_test

import (
	"testing"

	"github.com/go-zepto/zepto/plugins/linkeradmin"
	"github.com/go-zepto/zepto/plugins/linkeradmin/fields"
	"github.com/stretchr/testify/assert"
)

func TestResourceFieldEndpoint_AddField(t *testing.T) {
	r := linkeradmin.NewResource("Person")

	idField := fields.NewNumberField("id", nil)
	createdAtField := fields.NewDatetimeField("created_at", nil)
	updatedAtField := fields.NewDatetimeField("updated_at", nil)

	r.List.
		AddField(idField).
		AddField(createdAtField).
		AddField(updatedAtField)

	assert.Equal(t, []*fields.Field{
		&idField,
		&createdAtField,
		&updatedAtField,
	}, r.List.Fields)
}

func TestResourceFieldEndpoint_RemoveField(t *testing.T) {
	r := linkeradmin.NewResource("Person")

	idField := fields.NewNumberField("id", nil)
	createdAtField := fields.NewDatetimeField("created_at", nil)
	updatedAtField := fields.NewDatetimeField("updated_at", nil)

	r.List.Fields = []*fields.Field{
		&idField,
		&createdAtField,
		&updatedAtField,
	}

	r.List.RemoveField("created_at")

	assert.Equal(t, []*fields.Field{
		&idField,
		&updatedAtField,
	}, r.List.Fields)
}

func TestResourceFieldEndpoint_ReplaceField(t *testing.T) {
	r := linkeradmin.NewResource("Person")

	idField := fields.NewNumberField("id", nil)
	createdAtField := fields.NewDatetimeField("created_at", nil)
	updatedAtField := fields.NewDatetimeField("updated_at", nil)

	r.List.Fields = []*fields.Field{
		&idField,
		&createdAtField,
		&updatedAtField,
	}

	textCreateAtField := fields.NewTextField("created_at", nil)
	r.List.ReplaceField("created_at", textCreateAtField)

	assert.Equal(t, []*fields.Field{
		&idField,
		&textCreateAtField,
		&updatedAtField,
	}, r.List.Fields)
}
