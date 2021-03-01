package zepto

import "github.com/go-zepto/zepto/web"

type Plugin interface {
	Name() string
	Instance() interface{}
	PrependMiddlewares() []web.MiddlewareFunc
	AppendMiddlewares() []web.MiddlewareFunc
	OnCreated(z *Zepto)
	OnZeptoInit(z *Zepto)
	OnZeptoStart(z *Zepto)
	OnZeptoStop(z *Zepto)
}
