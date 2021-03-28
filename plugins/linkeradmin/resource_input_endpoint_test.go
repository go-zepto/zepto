package linkeradmin_test

import (
	"testing"

	"github.com/go-zepto/zepto/plugins/linkeradmin"
	"github.com/go-zepto/zepto/plugins/linkeradmin/fields"
	"github.com/stretchr/testify/assert"
)

func TestResourceInputEndpoint_AddInput(t *testing.T) {
	r := linkeradmin.NewResource("Person")

	idField := fields.NewNumberInput("id", nil)
	createdAtField := fields.NewDatetimeInput("created_at", nil)
	updatedAtField := fields.NewDatetimeInput("created_at", nil)

	r.Create.
		AddInput(idField).
		AddInput(createdAtField).
		AddInput(updatedAtField)

	assert.Equal(t, []*fields.Input{
		&idField,
		&createdAtField,
		&updatedAtField,
	}, r.Create.Inputs)
}

func TestResourceInputEndpoint_RemoveField(t *testing.T) {
	r := linkeradmin.NewResource("Person")

	idField := fields.NewNumberInput("id", nil)
	createdAtField := fields.NewDatetimeInput("created_at", nil)
	updatedAtField := fields.NewDatetimeInput("updated_at", nil)

	r.Create.Inputs = []*fields.Input{
		&idField,
		&createdAtField,
		&updatedAtField,
	}

	r.Create.RemoveInput("created_at")

	assert.Equal(t, []*fields.Input{
		&idField,
		&updatedAtField,
	}, r.Create.Inputs)
}

func TestResourceInputEndpoint_ReplaceField(t *testing.T) {
	r := linkeradmin.NewResource("Person")

	idField := fields.NewNumberInput("id", nil)
	createdAtField := fields.NewDatetimeInput("created_at", nil)
	updatedAtField := fields.NewDatetimeInput("created_at", nil)

	r.Create.Inputs = []*fields.Input{
		&idField,
		&createdAtField,
		&updatedAtField,
	}

	textCreatedAtInput := fields.NewTextInput("created_at", nil)
	r.Create.ReplaceInput("created_at", textCreatedAtInput)

	assert.Equal(t, []*fields.Input{
		&idField,
		&textCreatedAtInput,
		&updatedAtField,
	}, r.Create.Inputs)
}
