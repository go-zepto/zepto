package fields

import "github.com/go-zepto/zepto/plugins/linkeradmin"

type DatetimeFieldOptions struct {
	Label    string
	Sortable *bool
}

func NewDatetimeField(name string, opts *DatetimeFieldOptions) linkeradmin.Field {
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
		Type:    "datetime",
		Options: o,
	}
}

type DatetimeInputOptions struct {
	Label      string
	HelperText string
	FullWidth  *bool
	Disabled   *bool
}

func NewDatetimeInput(name string, opts *DatetimeInputOptions) linkeradmin.Input {
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
	}
	o["props"] = props
	return linkeradmin.Input{
		Name:    name,
		Type:    "datetime",
		Options: o,
	}
}
