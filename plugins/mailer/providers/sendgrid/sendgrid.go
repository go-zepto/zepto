package sendgrid

import (
	"errors"
	"net/http"
	"os"

	"github.com/go-zepto/zepto/plugins/mailer"
	"github.com/go-zepto/zepto/plugins/mailer/utils"
	"github.com/go-zepto/zepto/web/renderer"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendgridClient interface {
	Send(email *mail.SGMailV3) (*rest.Response, error)
}

type Settings struct {
	ApiKey      string
	DefaultFrom *mailer.Email
}

type SendgridProvider struct {
	rendererEngine renderer.Engine
	defaultFrom    *mailer.Email
	client         SendgridClient
}

func NewSendgridProvider(settings Settings) *SendgridProvider {
	apiKey := os.Getenv("SENDGRID_API_TOKEN")
	if apiKey == "" && settings.ApiKey == "" {
		panic("[mailer] sendgrid api key is required")
	}
	if apiKey == "" {
		apiKey = settings.ApiKey
	}
	return &SendgridProvider{
		defaultFrom: settings.DefaultFrom,
		client:      sendgrid.NewSendClient(apiKey),
	}
}

func (sg *SendgridProvider) Init(opts *mailer.InitOptions) {
	sg.rendererEngine = opts.RendererEngine
}

func (sg *SendgridProvider) SendFromTemplate(template string, opts *mailer.SendOptions) error {
	html, err := utils.RenderTemplateToString(sg.rendererEngine, template, opts.Vars)
	if err != nil {
		return err
	}
	return sg.SendFromHTML(html, opts)
}

func (sg *SendgridProvider) SendFromHTML(html string, opts *mailer.SendOptions) error {
	m := mail.NewV3Mail()

	var f = opts.From
	if f == nil {
		if sg.defaultFrom == nil {
			return errors.New("[mailer] missing From param")
		}
		f = sg.defaultFrom
	}

	from := mail.NewEmail(f.Name, f.Address)
	m.SetFrom(from)

	m.Subject = opts.Subject

	p := mail.NewPersonalization()
	tos := []*mail.Email{}
	for _, to := range opts.To {
		tos = append(tos, mail.NewEmail(to.Name, to.Address))
	}
	p.AddTos(tos...)

	m.AddPersonalizations(p)

	m.AddContent(mail.NewContent("text/html", html))
	if opts.PlainTextContent != "" {
		m.AddContent(mail.NewContent("text/plain", opts.PlainTextContent))
	}

	res, err := sg.client.Send(m)
	if err != nil {
		return err
	}
	if res.StatusCode > 299 {
		return errors.New("[sendgrid] " + http.StatusText(res.StatusCode))
	}
	return nil
}
