package fields

import "github.com/go-zepto/zepto/plugins/linkeradmin"

type CommonFieldOptions struct {
	Label     *string
	Sortable  *bool
	TextAlign *string
	EmptyText *string
}

func createLinkerAdminFieldOptions(c *CommonFieldOptions) linkeradmin.FieldOptions {
	o := make(linkeradmin.FieldOptions)
	o["props"] = map[string]interface{}{}
	if c.Label != nil {
		o["label"] = c.Label
	}
	if c.Sortable != nil {
		o["sortable"] = c.Sortable
	}
	if c.TextAlign != nil {
		o["textAlign"] = c.TextAlign
	}
	if c.EmptyText != nil {
		o["emptyText"] = c.EmptyText
	}
	return o
}
