package core

type Router struct {
	handlers map[any]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		handlers: map[any]HandlerFunc{},
	}
}

func (router *Router) AddHandler(label any, handler HandlerFunc) {
	router.handlers[label] = handler
}

func (router *Router) AddDefaultHandler(handler HandlerFunc) {
	router.handlers[nil] = handler
}

func (router *Router) Handler() HandlerFunc {
	return func(ctx *HandlerContext, err error) {
		defaultHandler, hasDefaultHandler := router.handlers[nil]
		handler, handlerExists := router.handlers[ctx.Options.Label]

		// If there's no label, or a handler for the current
		// label doesn't exist, fall back to the default handler.
		if ctx.Options.Label == nil || !handlerExists {
			// If no default handler, do absolutely nothing.
			if !hasDefaultHandler {
				return
			}

			defaultHandler(ctx, err)
			return
		}

		handler(ctx, err)
	}
}
