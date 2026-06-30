package goroute

type RouterGroup struct {
	prefix      string
	middlewares []RouteHandler //Middlewares for this group
	engine      *Engine        // Pointer to the main engine to access the router
}

// Group creates a new sub-group from the current group
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		engine: engine,
	}
	newGroup.middlewares = make([]RouteHandler, len(group.middlewares))
	copy(newGroup.middlewares, group.middlewares)
	return newGroup
}

// Use adds middlewares ONLY to the current group.
func (group *RouterGroup) Use(middlewares ...RouteHandler) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// addRoute pre-calculate the complete execution chain before the server even starts.
func (group *RouterGroup) addRoute(method string, comp string, handler RouteHandler) {
	pattern := group.prefix + comp

	handlers := make([]RouteHandler, len(group.middlewares), len(group.middlewares)+1)
	copy(handlers, group.middlewares)
	handlers = append(handlers, handler)
	group.engine.router.addRoute(method, pattern, handlers)
}

func (group *RouterGroup) GET(pattern string, handler RouteHandler) {
	group.addRoute("GET", pattern, handler)
}
func (group *RouterGroup) POST(pattern string, handler RouteHandler) {
	group.addRoute("POST", pattern, handler)
}
func (group *RouterGroup) PUT(pattern string, handler RouteHandler) {
	group.addRoute("PUT", pattern, handler)
}
func (group *RouterGroup) PATCH(pattern string, handler RouteHandler) {
	group.addRoute("PATCH", pattern, handler)
}
func (group *RouterGroup) DELETE(pattern string, handler RouteHandler) {
	group.addRoute("DELETE", pattern, handler)
}
