package web

import "github.com/go-zepto/zepto/mailer"

type Mailer struct {
	provider mailer.Mailer
}

func (m *Mailer) SendFromHTML(html string, opts *mailer.SendOptions) error {
	return m.provider.SendFromHTML(html, opts)
}

func (m *Mailer) SendFromTemplate(template string, opts *mailer.SendOptions) error {
	return m.provider.SendFromTemplate(template, opts)
}
