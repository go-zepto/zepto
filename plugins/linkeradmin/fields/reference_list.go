package fields

import "github.com/go-zepto/zepto/plugins/linkeradmin"

type ReferenceListFieldOptions struct {
	Label         string
	TextFieldName string
}

func NewReferenceListField(name string, resourceName string, opts *ReferenceListFieldOptions) linkeradmin.Field {
	o := make(linkeradmin.FieldOptions)
	props := map[string]interface{}{}
	if opts != nil {
		if opts.Label != "" {
			props["label"] = opts.Label
		}
		o["text_field_name"] = opts.TextFieldName
	}
	o["props"] = props
	o["ref_resource"] = resourceName
	o["ref_type"] = "list"
	return linkeradmin.Field{
		Name:    name,
		Type:    "reference_list",
		Options: o,
	}
}

type ReferenceListInputAutocomplete struct {
	Enabled          bool     `json:"enabled"`
	SearchableFields []string `json:"searchable_fields"`
}

type ReferenceListInputOptions struct {
	Label      string
	HelperText string
	Disabled   *bool
}

func NewReferenceListInput(resourceName string, resourceFieldName string, opts *ReferenceListInputOptions) linkeradmin.Input {
	o := make(linkeradmin.FieldOptions)
	props := map[string]interface{}{}
	if opts != nil {
		if opts.Label != "" {
			props["label"] = opts.Label
		}
	}
	o["props"] = props
	o["ref_resource"] = resourceName
	o["ref_resource_field_name"] = resourceFieldName
	o["ref_type"] = "list"
	return linkeradmin.Input{
		Name:    "reference_" + resourceFieldName,
		Type:    "reference_list",
		Options: o,
	}
}
