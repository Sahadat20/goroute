package goroute

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request
	Path   string
	Method string
	//value extracted from dynamic route
	Params map[string]string

	// The execution chain and state pointer
	handlers []RouteHandler
	index    int8                   //tracks which handler in the chain are currently executing
	Keys     map[string]interface{} // local storage for request-scoped context variables

}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
		index:  -1, // Start at -1, so the first Next() call increments it to 0
	}
}

// set attaches a key-value pair to the current request lifecycle context

func (c *Context) Set(key string, value interface{}) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = value
}

// Get retrieves a key-value pair from the context, returning false if missing
func (c *Context) Get(key string) (value interface{}, exists bool) {
	if c.Keys == nil {
		return nil, false
	}
	value, exists = c.Keys[key]
	return
}

// Next executes the pending handlers in the chain inside the calling handler.
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}

}

// Abort flags the framework to immediately stop driving the execution chain.
func (c *Context) Abort() {
	// Overwrite the index pointer to the maximum length to break the loop condition
	c.index = int8(len(c.handlers))
}

// Param retrieves a dynamic path parameter by its name
func (c *Context) Param(key string) string {
	return c.Params[key]
}

// String sends a plain text response with a status code
func (c *Context) String(code int, text string) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(code)
	c.Writer.Write([]byte(text))

}

// JSON sends a formatted JSON response with a status code
func (c *Context) JSON(code int, obj interface{}) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(code)

	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}

}

// BindJSON reads the incoming HTTP request body and decodes it into a Go struct.
func (c *Context) BindJSON(obj interface{}) error {
	decoder := json.NewDecoder(c.Req.Body)
	defer c.Req.Body.Close()
	return decoder.Decode(obj)
}

// SetHeader sets a specific header key-value pair in the HTTP response.
func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}
