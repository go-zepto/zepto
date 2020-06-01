package web

import (
	"context"
	"github.com/go-zepto/zepto/broker"
	"github.com/go-zepto/zepto/logger"
	"github.com/go-zepto/zepto/web/renderer"
	"github.com/gorilla/mux"
	"net/http"
	"sync"
	// Enable webpack asset feature
	_ "github.com/go-webpack/pongo2"
)

type Context interface {
	context.Context
	Vars() map[string]string
	Set(string, interface{})
	SetStatus(status int) Context
	Render(template string) error
	Logger() logger.Logger
	Broker() *broker.Broker
	Session() *Session
}

type DefaultContext struct {
	logger logger.Logger
	broker *broker.Broker
	context.Context
	res        http.ResponseWriter
	req        *http.Request
	status     int
	data       *sync.Map
	tmplEngine renderer.Engine
	session    *Session
}

func NewDefaultContext() *DefaultContext {
	return &DefaultContext{
		Context: context.Background(),
		data:    &sync.Map{},
		status:  200,
	}
}

func (d *DefaultContext) Set(key string, value interface{}) {
	d.data.Store(key, value)
}

func (d *DefaultContext) Value(key interface{}) interface{} {
	if k, ok := key.(string); ok {
		if v, ok := d.data.Load(k); ok {
			return v
		}
	}
	return d.Context.Value(key)
}

func (d *DefaultContext) SetStatus(s int) Context {
	d.status = s
	return d
}

func (d *DefaultContext) Render(template string) error {
	return d.tmplEngine.Render(d.res, d.status, template, d.data)
}

func (d *DefaultContext) Logger() logger.Logger {
	return d.logger
}

func (d *DefaultContext) Broker() *broker.Broker {
	return d.broker
}

func (d *DefaultContext) Vars() map[string]string {
	return mux.Vars(d.req)
}

func (d *DefaultContext) Session() *Session {
	return d.session
}
