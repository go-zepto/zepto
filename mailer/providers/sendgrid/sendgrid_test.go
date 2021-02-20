package sendgrid

import (
	"fmt"
	"testing"

	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/mailer"
)

func TestSendgrid(t *testing.T) {
	sg := NewSendgridProvider(Settings{
		ApiKey: "abc",
	})
	fmt.Println(sg)
}

func TestSendgridInit(t *testing.T) {
	z := zepto.NewZepto(zepto.Name("mailer-service"))
	sg := NewSendgridProvider(Settings{
		ApiKey: "abc",
	})
	sg.Init(mailer.InitOptions{
		Zepto: z,
	})
}
