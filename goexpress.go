package goroute

import (
	"fmt"
	"net/http"
)

type RouteHandler func(c *Context)

type Engine struct {
	router map[string]RouteHandler
}

func New() *Engine {
	return &Engine{router: make(map[string]RouteHandler)}
}

func (e *Engine) GET(path string, handler RouteHandler) {
	e.router["GET-"+path] = handler
}
func (e *Engine) POST(path string, handler RouteHandler) {
	e.router["POST-"+path] = handler
}
func (e *Engine) PUT(path string, handler RouteHandler) {
	e.router["PUT-"+path] = handler
}
func (e *Engine) DELETE(path string, handler RouteHandler) {
	e.router["DELETE-"+path] = handler
}
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Method + "-" + r.URL.Path
	if routeHandler, ok := e.router[key]; ok {
		c := newContext(w, r)
		routeHandler(c)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL)
	}
}
func (e *Engine) Run(addr string) error {
	fmt.Printf("GoExpress is running on %s...\n", addr)
	return http.ListenAndServe(addr, e)
}
