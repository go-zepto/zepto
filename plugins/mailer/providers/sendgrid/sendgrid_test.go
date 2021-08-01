package sendgrid

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/plugins/mailer"
	"github.com/go-zepto/zepto/web"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/assert"
)

type clientMock struct {
	statusCode int
	mail       *mail.SGMailV3
}

func (cm *clientMock) setStatus(status int) {
	cm.statusCode = status
}

func (cm *clientMock) Send(email *mail.SGMailV3) (*rest.Response, error) {
	cm.mail = email
	sc := cm.statusCode
	if sc == 0 {
		sc = 200
	}
	return &rest.Response{
		StatusCode: sc,
	}, nil
}

func setupZeptoWithMailMock() (*zepto.Zepto, *clientMock) {
	z := zepto.NewZepto()
	sg := NewSendgridProvider(Settings{
		ApiKey: "abc",
	})
	sg.Init(&mailer.InitOptions{
		RendererEngine: z.RendererEngine(),
	})
	cm := &clientMock{}
	sg.client = cm
	z.AddPlugin(mailer.NewMailerPlugin(mailer.Options{
		Mailer: sg,
	}))
	return z, cm
}

func TestSendgrid(t *testing.T) {
	setupZeptoWithMailMock()
}

var DEFAULT_SEND_OPTIONS = &mailer.SendOptions{
	Subject:          "Email in HTML",
	PlainTextContent: "Some email in HTML format",
	From: &mailer.Email{
		Name:    "From Name",
		Address: "from@test.com",
	},
	To: []*mailer.Email{
		{
			Name:    "To Name",
			Address: "to@test.com",
		},
	},
}

func TestSendgridSendMail(t *testing.T) {
	z, cm := setupZeptoWithMailMock()
	cm.setStatus(200)
	z.Post("/send-mail", func(ctx web.Context) error {
		m := mailer.InstanceFromCtx(ctx)
		err := m.SendFromHTML("<bold>Some email in HTML format</bold>", DEFAULT_SEND_OPTIONS)
		if err != nil {
			return err
		}
		return ctx.RenderJson(map[string]bool{
			"status": true,
		})
	})
	z.Init()
	z.InitApp()
	w := httptest.NewRecorder()
	z.App.ServeHTTP(w, httptest.NewRequest("POST", "/send-mail", nil))
	body := w.Body.String()
	assert.NotEmpty(t, body)
	m := cm.mail

	assert.Equal(t, 200, w.Result().StatusCode)
	assert.Equal(t, "Email in HTML", m.Subject)
	assert.Equal(t, "From Name", m.From.Name)
	assert.Equal(t, "from@test.com", m.From.Address)

	assert.Len(t, m.Personalizations, 1)
	p := m.Personalizations[0]
	assert.Len(t, p.To, 1)
	assert.Equal(t, "To Name", p.To[0].Name)
	assert.Equal(t, "to@test.com", p.To[0].Address)

	assert.Len(t, m.Content, 2)
	assert.Equal(t, "text/html", m.Content[0].Type)
	assert.Equal(t, "<bold>Some email in HTML format</bold>", m.Content[0].Value)
	assert.Equal(t, "text/plain", m.Content[1].Type)
	assert.Equal(t, "Some email in HTML format", m.Content[1].Value)
}

func TestSendgridSendMail_Error(t *testing.T) {
	os.Setenv("ZEPTO_ENV", "development")
	z, cm := setupZeptoWithMailMock()
	cm.setStatus(401)
	z.Post("/send-mail", func(ctx web.Context) error {
		m := mailer.InstanceFromCtx(ctx)
		err := m.SendFromHTML("<bold>Some email in HTML format</bold>", DEFAULT_SEND_OPTIONS)
		if err != nil {
			return err
		}
		return ctx.RenderJson(map[string]bool{
			"status": true,
		})
	})
	z.InitApp()
	w := httptest.NewRecorder()
	z.App.ServeHTTP(w, httptest.NewRequest("POST", "/send-mail", nil))
	body := w.Body.String()
	assert.NotEmpty(t, body)
	assert.Contains(t, body, "internal server error: [sendgrid] Unauthorized")
}
