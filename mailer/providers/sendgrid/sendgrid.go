package sendgrid

import (
	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/mailer"
)

type Settings struct {
	ApiKey string
}

type SendgridProvider struct {
	z      *zepto.Zepto
	apiKey string
}

func NewSendgridProvider(settings Settings) *SendgridProvider {
	return &SendgridProvider{
		apiKey: settings.ApiKey,
	}
}

func (sg *SendgridProvider) Init(opts mailer.InitOptions) {
	sg.z = opts.Zepto
}

func (sg *SendgridProvider) SendFromHTML(html string, opts *mailer.SendOptions) {

}
