package fields

type DatetimeFieldOptions struct {
	Label    string
	Sortable *bool
	ShowTime *bool
}

func NewDatetimeField(name string, opts *DatetimeFieldOptions) Field {
	o := make(FieldOptions)
	props := map[string]interface{}{}
	if opts != nil {
		if opts.Label != "" {
			props["label"] = opts.Label
		}
		if opts.Sortable != nil {
			props["sortable"] = opts.Sortable
		}
		if opts.ShowTime != nil {
			props["showTime"] = opts.ShowTime
		}
	}
	o["props"] = props
	return Field{
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

func NewDatetimeInput(name string, opts *DatetimeInputOptions) Input {
	o := make(FieldOptions)
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
	return Input{
		Name:    name,
		Type:    "datetime",
		Options: o,
	}
}
