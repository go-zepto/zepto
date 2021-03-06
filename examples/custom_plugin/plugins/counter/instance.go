package counter

import "github.com/go-zepto/zepto/web"

type Counter struct {
	value int
}

func (c *Counter) Value() int {
	return c.value
}

func (c *Counter) Add(n int) {
	c.value += n
}

func (c *Counter) Sub(n int) {
	c.value -= n
}

func (c *Counter) Reset() {
	c.value = 0
}

func InstanceFromCtx(ctx web.Context) *Counter {
	i := ctx.PluginInstance("counter")
	counter, ok := i.(*Counter)
	if !ok {
		return nil
	}
	return counter
}
