package fields

type TextFieldOptions struct {
	Label    string
	Sortable *bool
}

func NewTextField(name string, opts *TextFieldOptions) Field {
	o := make(FieldOptions)
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
	return Field{
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
	Disabled     *bool
}

func NewTextInput(name string, opts *TextInputOptions) Input {
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
		if opts.InitialValue != "" {
			props["initialValue"] = opts.InitialValue
		}
		if opts.Disabled != nil {
			props["disabled"] = opts.Disabled
		}
	}
	o["props"] = props
	return Input{
		Name:    name,
		Type:    "text",
		Options: o,
	}
}
