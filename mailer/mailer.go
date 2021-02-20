package mailer

import (
	"github.com/go-zepto/zepto"
)

type InitOptions struct {
	Zepto *zepto.Zepto
}

type Email struct {
	Name    string
	Address string
}

type SendOptions struct {
	Subject          string
	From             *Email
	To               []*Email
	PlainTextContent string
	Vars             map[string]interface{}
}

type Mailer interface {
	Init(opts *InitOptions)
	SendFromHTML(html string, opts *SendOptions) error
	// SendFromTemplate(template string, opts *SendOptions)
}
