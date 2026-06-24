package goroute

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// test 1: Does the Engine successfully find and execute a valid route?
func TestEngineRouting(t *testing.T) {
	// 1. setup a engine
	engine := New()

	// 2. Register a test route
	engine.GET("/hello", func(c *Context) {
		c.String(http.StatusOK, "hello world")
	})
	// 3. Create a fake request
	req := httptest.NewRequest("GET", "/hello", nil)

	// 4. create a fake responsewriter
	w := httptest.NewRecorder()

	// 5. execute engine directly
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
	if w.Body.String() != "hello world" {
		t.Errorf("Expected body 'hello world', got '%s'", w.Body.String())
	}
}

func TestEngineNotFound(t *testing.T) {
	engine := New()

	req := httptest.NewRequest("GET", "/missing-page", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code 404, got %d", w.Code)
	}
}

func TestContextJSON(t *testing.T) {
	engine := New()
	engine.GET("/api/data", func(c *Context) {
		c.JSON(http.StatusCreated, map[string]string{"status": "ok"})
	})

	// 3. Create a fake request
	req := httptest.NewRequest("GET", "/api/data", nil)

	// 4. create a fake responsewriter
	w := httptest.NewRecorder()

	// 5. execute engine directly
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code 201, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got '%s'", contentType)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}
	if response["status"] != "ok" {
		t.Errorf("Expected JSON status 'ok', got '%s'", response["status"])
	}

}

func TestPostAndBindJSON(t *testing.T) {
	engine := New()

	type Payload struct {
		Message string `json:"message"`
	}

	engine.POST("/echo", func(c *Context) {
		var p Payload
		if err := c.BindJSON(&p); err != nil {
			c.String(http.StatusBadRequest, "bad request")
			return
		}
		c.String(http.StatusCreated, "Received: "+p.Message)
	})

	// make request
	jsonBody := []byte(`{"message": "hello framework"}`)
	bodyReader := bytes.NewBuffer(jsonBody)
	req := httptest.NewRequest("POST", "/echo", bodyReader)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
	expectedBody := "Received: hello framework"
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, w.Body.String())
	}
}

// Test for the PUT method routing and JSON parsing
func TestPutRouting(t *testing.T) {
	engine := New()

	type UpdatePayload struct {
		Role string `json:"role"`
	}

	// 1. Register a PUT route for updates
	engine.PUT("/users/1/role", func(c *Context) {
		var p UpdatePayload
		if err := c.BindJSON(&p); err != nil {
			c.String(http.StatusBadRequest, "bad request")
			return
		}

		// If successful, confirm the update
		c.String(http.StatusOK, "Role updated to: "+p.Role)
	})

	// 2. Forge the incoming JSON Body
	jsonBody := []byte(`{"role": "admin"}`)
	bodyReader := bytes.NewBuffer(jsonBody)

	// 3. Create the fake PUT request
	req := httptest.NewRequest("PUT", "/users/1/role", bodyReader)
	w := httptest.NewRecorder()

	// 4. Execute
	engine.ServeHTTP(w, req)

	// 5. Assertions
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	expectedBody := "Role updated to: admin"
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, w.Body.String())
	}
}

// Test for the DELETE method routing
func TestDeleteRouting(t *testing.T) {
	engine := New()

	engine.DELETE("/remove", func(c *Context) {
		c.String(http.StatusOK, "deleted")
	})

	req := httptest.NewRequest("DELETE", "/remove", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestDynamicRoute(t *testing.T) {
	engine := New()

	// register a dynamic route
	engine.GET("/users/:id", func(c *Context) {
		id := c.Param("id")
		c.String(http.StatusOK, "User ID is "+id)
	})

	// test case 1: standard parameter extraction
	req := httptest.NewRequest("GET", "/users/456", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	expected := "User ID is 456"
	if w.Body.String() != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, w.Body.String())
	}
}
func TestWildcardRoute(t *testing.T) {
	engine := New()

	// Register a catch-all wildcard route
	engine.GET("/assets/*filepath", func(c *Context) {
		c.String(http.StatusOK, "File: "+c.Param("filepath"))
	})

	// Test Case 2: Multi-segment parameter capture
	req := httptest.NewRequest("GET", "/assets/css/main.css", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	expected := "File: css/main.css"
	if w.Body.String() != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, w.Body.String())
	}
}
func TestMiddlewareOrder(t *testing.T) {
	engine := New()

	// track the order of execution in this slice
	var executionOrder []string
	// Middleware A (Outer Layer)
	engine.Use(func(c *Context) {
		executionOrder = append(executionOrder, "A_Before")
		c.Next() // Suspend and go deeper
		executionOrder = append(executionOrder, "A_After")
	})

	// Middleware B (Inner Layer)
	engine.Use(func(c *Context) {
		executionOrder = append(executionOrder, "B_Before")
		c.Next() // Suspend and go to core handler
		executionOrder = append(executionOrder, "B_After")
	})

	// Core Route Handler
	engine.GET("/test", func(c *Context) {
		executionOrder = append(executionOrder, "Core_Handler")
		c.String(http.StatusOK, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	// Validate the Onion Model flow
	expectedOrder := []string{"A_Before", "B_Before", "Core_Handler", "B_After", "A_After"}

	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("Expected %d steps, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, v := range executionOrder {
		if v != expectedOrder[i] {
			t.Errorf("At index %d: expected %s, got %s", i, expectedOrder[i], v)
		}
	}
}
