package fields

import (
	"testing"

	"github.com/go-zepto/zepto/plugins/linkeradmin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/thriftrw/ptr"
)

func TestNewTextField_Defaults(t *testing.T) {
	f := NewTextField("name", nil)
	expectedField := linkeradmin.Field{
		Name: "name",
		Type: "text",
		Options: linkeradmin.FieldOptions{
			"props": map[string]interface{}{},
		},
	}
	assert.Equal(t, expectedField, f)
}

func TestNewTextField_WithOptions(t *testing.T) {
	f := NewTextField("name", &TextFieldOptions{
		Label:    "Name",
		Sortable: ptr.Bool(true),
	})
	expectedField := linkeradmin.Field{
		Name: "name",
		Type: "text",
		Options: linkeradmin.FieldOptions{
			"props": map[string]interface{}{
				"label":    "Name",
				"sortable": ptr.Bool(true),
			},
		},
	}
	assert.Equal(t, expectedField, f)
}

func TestNewTextInput_Defaults(t *testing.T) {
	f := NewTextInput("name", nil)
	expectedInput := linkeradmin.Input{
		Name: "name",
		Type: "text",
		Options: linkeradmin.FieldOptions{
			"props": map[string]interface{}{},
		},
	}
	assert.Equal(t, expectedInput, f)
}

func TestNewTextInput_WithOptions(t *testing.T) {
	f := NewTextInput("name", &TextInputOptions{
		Label:      "Name",
		HelperText: "The person name",
		FullWidth:  ptr.Bool(true),
		Disabled:   ptr.Bool(true),
	})
	expectedInput := linkeradmin.Input{
		Name: "name",
		Type: "text",
		Options: linkeradmin.InputOptions{
			"props": map[string]interface{}{
				"label":      "Name",
				"helperText": "The person name",
				"fullWidth":  ptr.Bool(true),
				"disabled":   ptr.Bool(true),
			},
		},
	}
	assert.Equal(t, expectedInput, f)
}
