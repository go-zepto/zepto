package mailer

import (
	"github.com/go-zepto/zepto/web/renderer"
)

type InitOptions struct {
	RendererEngine renderer.Engine
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
	SendFromTemplate(template string, opts *SendOptions) error
}
