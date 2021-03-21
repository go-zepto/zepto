package fields

import "github.com/go-zepto/zepto/plugins/linkeradmin"

type SelectFieldOptions struct {
	Label    string
	Sortable *bool
}

func NewSelectField(name string, opts *SelectFieldOptions) linkeradmin.Field {
	o := make(linkeradmin.FieldOptions)
	props := map[string]interface{}{}
	if opts != nil {
		if opts.Label != "" {
			props["label"] = opts.Label
		}
		if opts.Sortable != nil {
			props["sortable"] = opts.Sortable
		}
	}
	o["props"] = props
	return linkeradmin.Field{
		Name:    name,
		Type:    "select",
		Options: o,
	}
}

type SelectInputChoice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type SelectInputOptions struct {
	Label        string
	HelperSelect string
	FullWidth    *bool
	InitialValue string
	Disabled     *bool
	Choices      []SelectInputChoice
	EmptyText    string
}

func NewSelectInput(name string, opts *SelectInputOptions) linkeradmin.Input {
	o := make(linkeradmin.FieldOptions)
	props := map[string]interface{}{}
	if opts != nil {
		if opts.Label != "" {
			props["label"] = opts.Label
		}
		if opts.HelperSelect != "" {
			props["helperSelect"] = opts.HelperSelect
		}
		if opts.FullWidth != nil {
			props["fullWidth"] = opts.FullWidth
		}
		if opts.InitialValue != "" {
			props["initialValue"] = opts.InitialValue
		}
		if opts.Disabled != nil {
			props["disabled"] = opts.Disabled
		}
		if opts.Choices != nil {
			props["choices"] = opts.Choices
		}
		if opts.EmptyText != "" {
			props["emptyText"] = opts.EmptyText
		}
		props["optionText"] = "name"
		props["optionValue"] = "value"
	}
	o["props"] = props
	return linkeradmin.Input{
		Name:    name,
		Type:    "select",
		Options: o,
	}
}
