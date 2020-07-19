package web

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
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

type TestHandlerRequestOptions struct {
	Handler        RouteHandler
	Method         string
	Target         string
	Body           io.Reader
	InitialSession map[string]string
}

type TestHandlerRequestResult struct {
	t      *testing.T
	Ctx    *DefaultContext
	Status int
	Header http.Header
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

func (tr *TestHandlerRequestResult) AssertSessionValue(key string, value string) {
	assert.Equal(tr.t, value, tr.Ctx.session.Get(key))
}

func (tr *TestHandlerRequestResult) AssertHeaderValue(key string, value string) {
	assert.Equal(tr.t, value, tr.Header.Get(key))
}

func (zt *ZeptoTest) TestHandlerRequest(opts TestHandlerRequestOptions) (TestHandlerRequestResult, error) {
	req := httptest.NewRequest(opts.Method, opts.Target, opts.Body)
	res := httptest.NewRecorder()
	ctx := NewDefaultContext()
	ctx.session = zt.app.getSession(res, req)
	ctx.req = req
	ctx.res = res
	ctx.logger = log.New()
	ctx.tmplEngine = zt.app.tmplEngine
	for key, value := range opts.InitialSession {
		ctx.session.Set(key, value)
	}
	err := opts.Handler(ctx)
	if err != nil {
		return TestHandlerRequestResult{}, err
	}
	resp := res.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	return TestHandlerRequestResult{
		t:      zt.t,
		Ctx:    ctx,
		Status: resp.StatusCode,
		Header: resp.Header,
		Body:   string(body),
	}, nil
}
