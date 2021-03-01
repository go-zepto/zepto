package mailernotifier

import (
	"github.com/go-zepto/zepto/mailer"
	"github.com/go-zepto/zepto/plugins/auth/authcore"
)

type TemplateOptions struct {
	TemplatePath string
	Subject      string
}

type Templates struct {
	ResetPasswordToken *TemplateOptions
	PasswordReseted    *TemplateOptions
}

type Options struct {
	Mailer    mailer.Mailer
	Templates *Templates
}

type MailerNotifier struct {
	mailer    mailer.Mailer
	templates *Templates
}

func NewMailerNotifier(opts Options) *MailerNotifier {
	if opts.Templates == nil {
		opts.Templates = &Templates{}
	}

	if opts.Templates.ResetPasswordToken == nil {
		opts.Templates.ResetPasswordToken = &TemplateOptions{
			TemplatePath: "mailer/reset_password_token",
			Subject:      "Reset your password",
		}
	}

	if opts.Templates.PasswordReseted == nil {
		opts.Templates.PasswordReseted = &TemplateOptions{
			TemplatePath: "mailer/password_reseted",
			Subject:      "Password reset confirmation",
		}
	}

	return &MailerNotifier{
		mailer:    opts.Mailer,
		templates: opts.Templates,
	}
}

func (mn *MailerNotifier) NotifyResetPasswordToken(email string, token *authcore.Token, pid authcore.PID) error {
	return mn.mailer.SendFromTemplate(mn.templates.ResetPasswordToken.TemplatePath, &mailer.SendOptions{
		Subject: mn.templates.ResetPasswordToken.Subject,
		To: []*mailer.Email{
			{
				Address: email,
			},
		},
		Vars: map[string]interface{}{
			"token": token,
			"pid":   pid,
		},
	})
}

func (mn *MailerNotifier) NotifyPasswordReseted(email string, pid authcore.PID) error {
	panic("not implemented")
}
