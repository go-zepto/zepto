package web

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

type ZeptoTest struct {
	t   *testing.T
	app *App
}

func NewZeptoTest(t *testing.T, app *App) *ZeptoTest {
	return &ZeptoTest{
		t:   t,
		app: app,
	}
}

type TestHandlerRequestResult struct {
	t      *testing.T
	Status int
	Body   string
}

func (tr *TestHandlerRequestResult) AssertStatusCode(status int) {
	assert.Equal(tr.t, status, tr.Status)
}

func (tr *TestHandlerRequestResult) AssertBodyEquals(str string) {
	assert.Equal(tr.t, str, tr.Body)
}

func (tr *TestHandlerRequestResult) AssertBodyContains(str string) {
	assert.Contains(tr.t, tr.Body, str)
}

func (zt *ZeptoTest) TestHandlerRequest(h RouteHandler, m string, t string, b io.Reader) (TestHandlerRequestResult, error) {
	req := httptest.NewRequest(m, t, b)
	res := httptest.NewRecorder()
	ctx := NewDefaultContext()
	ctx.req = req
	ctx.res = res
	ctx.logger = log.New()
	ctx.tmplEngine = zt.app.tmplEngine
	err := h(ctx)
	if err != nil {
		return TestHandlerRequestResult{}, err
	}
	resp := res.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	return TestHandlerRequestResult{
		t:      zt.t,
		Status: resp.StatusCode,
		Body:   string(body),
	}, nil
}
