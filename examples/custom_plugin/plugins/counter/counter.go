package counter

import (
	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/web"
)

/*
	CounterPlugin adds 5 routes to control a stateful counter.
	- GET /counter (View the counter value)
	- GET /counter/inc (Increment 1 to the counter)
	- GET /counter/dec (Decrement 1 to the counter)
	- GET /counter/reset (Set counter value to zero)

	The instance is the Counter object and can be easily accessed in web.Context
*/

type CounterPlugin struct {
	counter *Counter
}

func NewCounterPlugin() *CounterPlugin {
	return &CounterPlugin{
		counter: &Counter{},
	}
}

func (cp *CounterPlugin) Name() string {
	return "counter"
}

func (cp *CounterPlugin) Instance() interface{} {
	return cp.counter
}

func (cp *CounterPlugin) PrependMiddlewares() []web.MiddlewareFunc {
	return nil
}

func (cp *CounterPlugin) AppendMiddlewares() []web.MiddlewareFunc {
	return nil
}

func (cp *CounterPlugin) OnCreated(z *zepto.Zepto) {
}

func (cp *CounterPlugin) OnZeptoInit(z *zepto.Zepto) {
	z.Get("/counter", func(ctx web.Context) error {
		return ctx.RenderJson(map[string]int{
			"count": cp.counter.Value(),
		})
	})
	z.Get("/counter/inc", func(ctx web.Context) error {
		cp.counter.Add(1)
		return ctx.RenderJson(map[string]int{
			"count": cp.counter.Value(),
		})
	})
	z.Get("/counter/dec", func(ctx web.Context) error {
		cp.counter.Sub(1)
		return ctx.RenderJson(map[string]int{
			"count": cp.counter.Value(),
		})
	})
	z.Get("/counter/reset", func(ctx web.Context) error {
		cp.counter.Reset()
		return ctx.RenderJson(map[string]int{
			"count": cp.counter.Value(),
		})
	})
}

func (cp *CounterPlugin) OnZeptoStart(z *zepto.Zepto) {

}

func (cp *CounterPlugin) OnZeptoStop(z *zepto.Zepto) {

}
