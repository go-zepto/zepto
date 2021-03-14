package fields

import "github.com/go-zepto/zepto/plugins/linkeradmin"

type ReferenceFieldOptions struct {
	Label         string
	Sortable      *bool
	TextFieldName string
}

func NewReferenceField(name string, resourceName string, opts *ReferenceFieldOptions) linkeradmin.Field {
	o := make(linkeradmin.FieldOptions)
	props := map[string]interface{}{}
	if opts != nil {
		if opts.Label != "" {
			props["label"] = opts.Label
		}
		if opts.Sortable != nil {
			props["sortable"] = opts.Sortable
		}
		o["text_field_name"] = opts.TextFieldName
	}
	o["props"] = props
	o["ref_resource"] = resourceName
	o["ref_type"] = "single"
	return linkeradmin.Field{
		Name:    name,
		Type:    "reference",
		Options: o,
	}
}

type ReferenceInputAutocomplete struct {
	Enabled          bool     `json:"enabled"`
	SearchableFields []string `json:"searchable_fields"`
}

type ReferenceInputOptions struct {
	Label           string
	HelperText      string
	FullWidth       *bool
	InitialValue    string
	Disabled        *bool
	OptionTextField string
	Autocomplete    ReferenceInputAutocomplete
}

func NewReferenceInput(name string, resourceName string, opts *ReferenceInputOptions) linkeradmin.Input {
	o := make(linkeradmin.FieldOptions)
	props := map[string]interface{}{}
	if opts != nil {
		if opts.Label != "" {
			props["label"] = opts.Label
		}
		if opts.Autocomplete.SearchableFields == nil {
			opts.Autocomplete.SearchableFields = make([]string, 0)
		}
		o["autocomplete"] = opts.Autocomplete
	}
	o["props"] = props
	o["ref_resource"] = resourceName
	o["ref_type"] = "single"
	o["option_text_field"] = opts.OptionTextField
	return linkeradmin.Input{
		Name:    name,
		Type:    "reference",
		Options: o,
	}
}
