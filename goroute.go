package goroute

import (
	"fmt"
	"net/http"
)

type RouteHandler func(c *Context)

type Engine struct {
	*RouterGroup
	router *router
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{
		engine: engine,
	}
	return engine
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// 1. Create the briefcase
	c := newContext(w, r)

	// 2. Pre-load the global middlewares into the execution chain
	c.handlers = append(c.handlers, e.middlewares...)

	// 3. find a specific route logic
	node, params := e.router.getRoute(r.Method, r.URL.Path)

	if node != nil {
		// 4. Inject the extracted dynamic parameters into the context
		c.Params = params

		// 3. append the core handler to the end of the chain
		key := r.Method + "-" + node.pattern
		c.handlers = e.router.handlers[key]
	} else {
		// Append a 404 handler to the chain if route is missing
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND")
		})
	}

	// 5. Kick off the execution chain
	c.Next()
}
func (e *Engine) Run(addr string) error {
	fmt.Printf("GoExpress is running on %s...\n", addr)
	return http.ListenAndServe(addr, e)
}
