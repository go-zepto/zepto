package mailerstub

import "github.com/go-zepto/zepto/mailer"

type SendFromHtmlStackItem struct {
	Html string
	Opts *mailer.SendOptions
}

type SendFromTemplateStackItem struct {
	Template string
	Opts     *mailer.SendOptions
}

type MailerStub struct {
	sendFromHtmlStack     []SendFromHtmlStackItem
	sendFromTemplateStack []SendFromTemplateStackItem
}

func NewMailerStub() *MailerStub {
	return &MailerStub{
		sendFromHtmlStack:     make([]SendFromHtmlStackItem, 0),
		sendFromTemplateStack: make([]SendFromTemplateStackItem, 0),
	}
}

func (ms *MailerStub) Init(opts *mailer.InitOptions) {}

func (ms *MailerStub) SendFromHTML(html string, opts *mailer.SendOptions) error {
	ms.sendFromHtmlStack = append(ms.sendFromHtmlStack, SendFromHtmlStackItem{
		Html: html,
		Opts: opts,
	})
	return nil
}

func (ms *MailerStub) SendFromTemplate(template string, opts *mailer.SendOptions) error {
	ms.sendFromTemplateStack = append(ms.sendFromTemplateStack, SendFromTemplateStackItem{
		Template: template,
		Opts:     opts,
	})
	return nil
}

func (ms *MailerStub) SentFromHTMLAt(idx int) *SendFromHtmlStackItem {
	return &ms.sendFromHtmlStack[idx]
}

func (ms *MailerStub) LastSentFromHTML() *SendFromHtmlStackItem {
	return &ms.sendFromHtmlStack[len(ms.sendFromHtmlStack)-1]
}

func (ms *MailerStub) SentFromTemplateAt(idx int) *SendFromTemplateStackItem {
	return &ms.sendFromTemplateStack[idx]
}

func (ms *MailerStub) LastSentFromTemplate() *SendFromTemplateStackItem {
	return &ms.sendFromTemplateStack[len(ms.sendFromTemplateStack)-1]
}
