package goroute

import (
	"log"
	"time"
)

func Logger() RouteHandler {
	return func(c *Context) {
		// start the timer
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		// 4. Print the structured log with formatted spacing
		log.Printf("[GoExpress] | %3d | %13v | %-7s %s",
			c.StatusCode, // Captured from our context.go upgrade
			latency,      // Total CPU time spent
			c.Method,
			c.Path,
		)
	}
}
