package web

import (
	"context"
	"github.com/go-zepto/zepto/broker"
	"github.com/go-zepto/zepto/web/renderer"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
	// Enable webpack asset feature
	_ "github.com/go-webpack/pongo2"
)

type Context interface {
	context.Context
	Set(string, interface{})
	SetStatus(status int) Context
	Render(template string) error
	Logger() *log.Logger
	Broker() *broker.Broker
}

type DefaultContext struct {
	logger *log.Logger
	broker *broker.Broker
	context.Context
	res        http.ResponseWriter
	req        *http.Request
	status     int
	data       *sync.Map
	tmplEngine renderer.Engine
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

func (d *DefaultContext) Logger() *log.Logger {
	return d.logger
}

func (d *DefaultContext) Broker() *broker.Broker {
	return d.broker
}
