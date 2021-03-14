package fields

import "github.com/go-zepto/zepto/plugins/linkeradmin"

type TextFieldOptions struct {
	Label    string
	Sortable *bool
}

func NewTextField(name string, opts *TextFieldOptions) linkeradmin.Field {
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
		Type:    "text",
		Options: o,
	}
}

type TextInputOptions struct {
	Label        string
	HelperText   string
	FullWidth    *bool
	InitialValue string
}

func NewTextInput(name string, opts *TextInputOptions) linkeradmin.Input {
	o := make(linkeradmin.FieldOptions)
	props := map[string]interface{}{}
	if opts != nil {
		if opts.Label != "" {
			props["label"] = opts.Label
		}
		if opts.HelperText != "" {
			props["helperText"] = opts.HelperText
		}
		if opts.FullWidth != nil {
			props["fullWidth"] = opts.FullWidth
		}
		if opts.InitialValue != "" {
			props["initialValue"] = opts.InitialValue
		}
	}
	o["props"] = props
	return linkeradmin.Input{
		Name:    name,
		Type:    "text",
		Options: o,
	}
}
