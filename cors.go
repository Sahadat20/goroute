package goroute

import "net/http"

func CORS() RouteHandler {
	return func(c *Context) {
		// 1. Inject the necessary CORS Headers into EVERY response
		c.SetHeader("Access-Control-Allow-Origin", "*")
		c.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// Whitelist the HTTP methods framework allows
		c.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.SetHeader("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Admin-Token")

		// 2. Handle the Preflight Request
		if c.Method == "OPTIONS" {
			// A preflight request doesn't need a response body.
			// It just needs the headers (which we attached above) and a 204 status.
			c.Writer.WriteHeader(http.StatusNoContent)

			// CRITICAL: We use Abort() from Lab 6 so the execution chain stops here!
			// We do not want preflight requests bleeding into our core logic.
			c.Abort()
			return
		}
		// 3. If it's a normal request (GET, POST, etc.), continue down the chain
		c.Next()
	}
}
