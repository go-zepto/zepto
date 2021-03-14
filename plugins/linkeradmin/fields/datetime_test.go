package fields

import (
	"testing"

	"github.com/go-zepto/zepto/plugins/linkeradmin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/thriftrw/ptr"
)

func TestNewDatetimeField_Defaults(t *testing.T) {
	f := NewDatetimeField("created_at", nil)
	expectedField := linkeradmin.Field{
		Name: "created_at",
		Type: "datetime",
		Options: linkeradmin.FieldOptions{
			"props": map[string]interface{}{},
		},
	}
	assert.Equal(t, expectedField, f)
}

func TestNewDatetimeField_WithOptions(t *testing.T) {
	f := NewDatetimeField("created_at", &DatetimeFieldOptions{
		Label:    "Created At",
		Sortable: ptr.Bool(true),
	})
	expectedField := linkeradmin.Field{
		Name: "created_at",
		Type: "datetime",
		Options: linkeradmin.FieldOptions{
			"props": map[string]interface{}{
				"label":    "Created At",
				"sortable": ptr.Bool(true),
			},
		},
	}
	assert.Equal(t, expectedField, f)
}

func TestNewDatetimeInput_Defaults(t *testing.T) {
	f := NewDatetimeInput("birthdate", nil)
	expectedInput := linkeradmin.Input{
		Name: "birthdate",
		Type: "datetime",
		Options: linkeradmin.FieldOptions{
			"props": map[string]interface{}{},
		},
	}
	assert.Equal(t, expectedInput, f)
}

func TestNewDatetimeInput_WithOptions(t *testing.T) {
	f := NewDatetimeInput("birthdate", &DatetimeInputOptions{
		Label:      "Birthdate",
		HelperText: "The date of birth",
		FullWidth:  ptr.Bool(true),
		Disabled:   ptr.Bool(true),
	})
	expectedInput := linkeradmin.Input{
		Name: "birthdate",
		Type: "datetime",
		Options: linkeradmin.InputOptions{
			"props": map[string]interface{}{
				"label":      "Birthdate",
				"helperText": "The date of birth",
				"fullWidth":  ptr.Bool(true),
				"disabled":   ptr.Bool(true),
			},
		},
	}
	assert.Equal(t, expectedInput, f)
}
