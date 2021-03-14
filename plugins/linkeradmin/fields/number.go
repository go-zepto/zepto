package fields

import "github.com/go-zepto/zepto/plugins/linkeradmin"

type NumberFieldOptions struct {
	Label    string
	Sortable *bool
}

func NewNumberField(name string, opts *NumberFieldOptions) linkeradmin.Field {
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
		Type:    "number",
		Options: o,
	}
}

type NumberInputOptions struct {
	Label      string
	HelperText string
	FullWidth  *bool
	Disabled   *bool
	Max        *int64
	Min        *int64
	Step       *int64
}

func NewNumberInput(name string, opts *NumberInputOptions) linkeradmin.Input {
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
		if opts.Disabled != nil {
			props["disabled"] = opts.Disabled
		}
		if opts.Max != nil {
			props["max"] = opts.Max
		}
		if opts.Min != nil {
			props["min"] = opts.Min
		}
		if opts.Step != nil {
			props["step"] = opts.Step
		}
	}
	o["props"] = props
	return linkeradmin.Input{
		Name:    name,
		Type:    "number",
		Options: o,
	}
}
