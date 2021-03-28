package fields

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/thriftrw/ptr"
)

func TestNewNumberField_Defaults(t *testing.T) {
	f := NewNumberField("age", nil)
	expectedField := Field{
		Name: "age",
		Type: "number",
		Options: FieldOptions{
			"props": map[string]interface{}{},
		},
	}
	assert.Equal(t, expectedField, f)
}

func TestNewNumberField_WithOptions(t *testing.T) {
	f := NewNumberField("age", &NumberFieldOptions{
		Label:    "Name",
		Sortable: ptr.Bool(true),
	})
	expectedField := Field{
		Name: "age",
		Type: "number",
		Options: FieldOptions{
			"props": map[string]interface{}{
				"label":    "Name",
				"sortable": ptr.Bool(true),
			},
		},
	}
	assert.Equal(t, expectedField, f)
}

func TestNewNumberInput_Defaults(t *testing.T) {
	f := NewNumberInput("age", nil)
	expectedInput := Input{
		Name: "age",
		Type: "number",
		Options: FieldOptions{
			"props": map[string]interface{}{},
		},
	}
	assert.Equal(t, expectedInput, f)
}

func TestNewNumberInput_WithOptions(t *testing.T) {
	f := NewNumberInput("age", &NumberInputOptions{
		Label:      "Age",
		HelperText: "The person age",
		FullWidth:  ptr.Bool(true),
		Disabled:   ptr.Bool(true),
		Min:        ptr.Int64(0),
		Max:        ptr.Int64(120),
		Step:       ptr.Int64(1),
	})
	expectedInput := Input{
		Name: "age",
		Type: "number",
		Options: InputOptions{
			"props": map[string]interface{}{
				"label":      "Age",
				"helperText": "The person age",
				"fullWidth":  ptr.Bool(true),
				"disabled":   ptr.Bool(true),
				"min":        ptr.Int64(0),
				"max":        ptr.Int64(120),
				"step":       ptr.Int64(1),
			},
		},
	}
	assert.Equal(t, expectedInput, f)
}
