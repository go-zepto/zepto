package mailer

import "github.com/go-zepto/zepto/web"

type MailerInstance interface {
	SendFromHTML(html string, opts *SendOptions) error
	SendFromTemplate(template string, opts *SendOptions) error
}

func InstanceFromCtx(ctx web.Context) MailerInstance {
	i := ctx.PluginInstance("mailer")
	mailerInstance, ok := i.(MailerInstance)
	if !ok {
		return nil
	}
	return mailerInstance
}
