package counter

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
