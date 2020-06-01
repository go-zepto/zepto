package web

import (
	"github.com/go-zepto/zepto/testutils"
	"github.com/go-zepto/zepto/web/renderer"
	log "github.com/sirupsen/logrus"
	"sync"
	"testing"
)

type CreateTestAppOptions struct {
	tmplEngine renderer.Engine
}

func mapLen(m *sync.Map) int {
	length := 0
	m.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	return length
}

func CreateTestApp(t *testing.T, opts CreateTestAppOptions) *DefaultContext {
	res := testutils.NewResponseMock(t)
	ctx := NewDefaultContext()
	ctx.logger = log.New()
	ctx.res = res
	ctx.tmplEngine = opts.tmplEngine
	return ctx
}

func TestDefaultContext_SetStatus(t *testing.T) {
	te := &testutils.EngineMock{}
	c := CreateTestApp(t, CreateTestAppOptions{
		tmplEngine: te,
	})
	if c.status != 200 {
		t.Error("Initial status from context should be 200")
	}
	c.SetStatus(400)
	if c.status != 400 {
		t.Error("SetStatus(400) should change status to 400")
	}
}

func TestDefaultContext_Set_And_Value(t *testing.T) {
	te := &testutils.EngineMock{}
	c := CreateTestApp(t, CreateTestAppOptions{
		tmplEngine: te,
	})
	if mapLen(c.data) != 0 {
		t.Error("Initial context data should be empty")
	}
	c.Set("hello", "world")
	value := c.Value("hello")
	if value == nil {
		t.Error("c.Set(\"hello\") should add a new item to data")
	}
	if value != "world" {
		t.Error("map value should be equal \"world\"")
	}
}

func TestDefaultContext_Render(t *testing.T) {
	te := &testutils.EngineMock{}
	c := CreateTestApp(t, CreateTestAppOptions{
		tmplEngine: te,
	})
	if te.RenderCalled {
		t.Error("Render should not be called at this point")
	}
	err := c.Render("teste")
	if err != nil {
		t.Error(err)
	}
	if !te.RenderCalled {
		t.Error("Render should be called at this point")
	}
}
