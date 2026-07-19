package main

import (
	"fmt"
	"net/http"
	"time"

	goroute "github.com/Sahadat20/goroute"
)

const SecretKey = "this-is-test-secret"

// Middleware for /admin Group
func AdminGuard() goroute.RouteHandler {
	return func(c *goroute.Context) {
		token := c.Req.Header.Get("X-Admin-Token")
		if token != "super-secret-admin-pass" {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"status": "Rejected",
				"reason": "Administrative clearance token missing or invalid.",
			})
			c.Abort() // 🔥 Overwrites c.index to 3, stopping the loop machine
			return    // Short-circuit the execution chain!
		}
		fmt.Println("[GUARD] Clearance confirmed. Transitioning inward...")
		c.Next()
	}
}

// Middleware for /something-else Group
func ContextualTracker() goroute.RouteHandler {
	return func(c *goroute.Context) {
		fmt.Println("[TRACKER] Request routed directly into the 'Something-Else' cluster.")
		c.Next()
	}
}

func main() {
	g := goroute.New()

	// attach recovry first
	g.Use(goroute.Recovery())
	// Apply Global Middleware
	g.Use(goroute.Logger())
	g.Use(goroute.CORS()) //Inject CORS middleware globally

	// THE WILDCARD PREFLIGHT CATCHER
	// Any OPTIONS request will match this wildcard. The request will flow through:
	// [GlobalLogger -> CORS -> Empty Handler]
	// But because CORS calls c.Abort(), the empty handler is safely ignored!
	g.OPTIONS("/*cors", func(c *goroute.Context) {
		// Intentionally left blank. c.Abort() in the CORS middleware stops execution before reaching here.
	})

	g.GET("/dangerous", func(c *goroute.Context) {
		var pointer *int // 1. Allocates a pointer, defaults to nil (0x0)
		*pointer = 10    // 2. Attempts to write to memory address 0x0
	})
	// ==========================================
	// GROUP 1: /admin
	// ==========================================
	// UNPROTECTED ENDPOINT: Token Generation Factory
	g.POST("/login", func(c *goroute.Context) {
		// Mock login validation assuming user passed parameters successfully
		token, err := goroute.GenerateJWT("pro_programmer", SecretKey, 15*time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to forge token"})
			return
		}
		c.JSON(http.StatusOK, map[string]string{"token": token})
	})
	admin := g.Group("/admin")
	admin.Use(goroute.JWTAuth(SecretKey))
	{
		// Resolves to: PUT /admin/users/:id
		admin.PUT("/users/:id", func(c *goroute.Context) {
			id := c.Param("id")
			c.JSON(http.StatusOK, map[string]string{
				"action":  "PUT",
				"message": fmt.Sprintf("User structure for ID %s has been completely overwritten.", id),
			})
		})

		// Resolves to: DELETE /admin/users/:id
		admin.DELETE("/users/:id", func(c *goroute.Context) {
			id := c.Param("id")
			c.JSON(http.StatusOK, map[string]string{
				"action":  "DELETE",
				"message": fmt.Sprintf("User record with ID %s purged permanently from system storage.", id),
			})
		})
	}

	// ==========================================
	// GROUP 2: /something-else
	// ==========================================
	somethingElse := g.Group("/something-else")
	somethingElse.Use(ContextualTracker())
	{
		// Resolves to: PATCH /something-else/configurations
		somethingElse.PATCH("/configurations", func(c *goroute.Context) {
			c.JSON(http.StatusOK, map[string]string{
				"action":  "PATCH",
				"message": "Specific configuration delta settings applied successfully.",
			})
		})
	}

	// ProductPayload defines the expected JSON structure and its machine rules
	type ProductPayload struct {
		Name  string  `json:"name" validate:"required"`
		Price float64 `json:"price" validate:"required"`
		SKU   string  `json:"sku"` // Optional
	}
	g.POST("/api/products", func(c *goroute.Context) {
		var payload ProductPayload

		// 1. One line to map the memory and validate the structural rules
		if err := c.BindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		// 2. If we reach here, the struct is fully populated and 100% safe to use
		c.JSON(http.StatusCreated, map[string]interface{}{
			"message": "Product safely bound to memory",
			"product": payload.Name,
		})
	})
	g.Run(":8080")
}
