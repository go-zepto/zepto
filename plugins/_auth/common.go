package auth

import (
	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/plugins/auth/authcore"
	"github.com/go-zepto/zepto/plugins/auth/encoders/uuid"
	mailernotifier "github.com/go-zepto/zepto/plugins/auth/notifiers/mailer"
	"github.com/go-zepto/zepto/plugins/auth/stores/inmemory"
)

func CreateAuthCore(z *zepto.Zepto, opts authcore.AuthCore) *authcore.AuthCore {
	ac := authcore.AuthCore{}
	if opts.TokenEncoder == nil {
		opts.TokenEncoder = uuid.NewUUIDTokenEncoder()
	}
	if opts.Store == nil {
		opts.Store = inmemory.NewInMemoryStore()
	}
	if opts.Notifier == nil {
		opts.Notifier = mailernotifier.NewMailerNotifier(mailernotifier.Options{
			Mailer: z.Mailer(),
		})
	}
	ac.DS = opts.DS
	ac.TokenEncoder = opts.TokenEncoder
	ac.Store = opts.Store
	ac.Notifier = opts.Notifier
	return &ac
}
