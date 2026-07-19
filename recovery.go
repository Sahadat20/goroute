package goroute

import (
	"log"
	"net/http"
	"runtime"
)

func Recovery() RouteHandler {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 1024)
				n := runtime.Stack(buf, false)

				log.Printf("[RECOVERY] Panic recovered: \n%v\n%s", err, buf[:n])
				c.String(http.StatusInternalServerError, "500 Internal Server Error")
				c.Abort()
			}
		}()
		c.Next()
	}

}
