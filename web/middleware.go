package web

type MiddlewareStack struct {
	stack []MiddlewareFunc
	skips map[string]bool
}

func (ms *MiddlewareStack) Use(mw ...MiddlewareFunc) {
	ms.stack = append(ms.stack, mw...)
}

func (ms *MiddlewareStack) UsePrepend(mw ...MiddlewareFunc) {
	ms.stack = append(mw, ms.stack...)
}

func (ms *MiddlewareStack) handle(handler RouteHandler) RouteHandler {
	if len(ms.stack) == 0 {
		return handler
	}
	for idx := len(ms.stack) - 1; idx >= 0; idx-- {
		handler = ms.stack[idx](handler)
	}
	return handler
}
