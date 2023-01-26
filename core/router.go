package core

import (
	"fmt"
	"sync"
)

var DefaultDefaultHandler HandlerFunc = func(ctx *HandlerContext, err error) {
	if ctx.Response == nil {
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("An error occurred.")
	}

	fmt.Printf("%s responded with status %s\n", ctx.Response.Request.URL, ctx.Response.Status)
}

type Router struct {
	handlers       map[any][]HandlerFunc
	defaultHandler HandlerFunc
	parallel       bool
}

// Creates a new router. Pass in "false" to run multiple handlers for a specific
// request sequentially. Pass "true" to run them in parallel
func NewRouter(parallel bool) *Router {
	return &Router{
		handlers:       map[any][]HandlerFunc{},
		defaultHandler: DefaultDefaultHandler,
		parallel:       parallel,
	}
}

func (router *Router) AddHandler(label any, handler HandlerFunc) {
	router.handlers[label] = append(router.handlers[label], handler)
}

func (router *Router) AddDefaultHandler(handler HandlerFunc) {
	router.defaultHandler = handler
}

func (router *Router) Handler() HandlerFunc {
	return func(ctx *HandlerContext, err error) {
		// Handle no label
		if ctx.Options.Label == nil {
			// If there is no default handler, do nothing.
			if router.defaultHandler == nil {
				return
			}

			// Otherwise, call the default handler.
			router.defaultHandler(ctx, err)
			return
		}

		handlers, handlersExist := router.handlers[ctx.Options.Label]

		if !handlersExist {
			return
		}

		if len(handlers) == 1 {
			handlers[0](ctx, err)
			return
		}

		if !router.parallel {
			for _, handler := range handlers {
				handler(ctx, err)
			}
			return
		}

		wg := sync.WaitGroup{}
		wg.Add(len(handlers))
		for _, handler := range handlers {
			go func(handler HandlerFunc) {
				defer wg.Done()
				handler(ctx, err)
			}(handler)
		}
		wg.Wait()
	}
}
