package fields

type ReferenceFieldOptions struct {
	Label         string
	Sortable      *bool
	TextFieldName string
}

func NewReferenceField(name string, resourceName string, opts *ReferenceFieldOptions) Field {
	o := make(FieldOptions)
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
	return Field{
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
	Filter          map[string]interface{}
}

func NewReferenceInput(name string, resourceName string, opts *ReferenceInputOptions) Input {
	o := make(FieldOptions)
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
	o["filter"] = opts.Filter
	return Input{
		Name:    name,
		Type:    "reference",
		Options: o,
	}
}
