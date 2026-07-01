package goroute

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSMidleware(t *testing.T) {
	// 1. setup a engine
	engine := New()
	engine.Use(CORS())

	var coreHandlerExecuted bool
	engine.GET("/test", func(c *Context) {
		coreHandlerExecuted = true
		c.String(http.StatusOK, "CORS test passed")
	})

	// ---Test 1: standard GET request---
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	// verify headers were injected
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("Expected Access-Control-Allow-Origin header to be '*', got '%s'", w.Header().Get("Access-Control-Allow-Origin"))
	}

	// test 2: Preflight OPTIONS request
	coreHandlerExecuted = false // Reset state
	req2 := httptest.NewRequest("OPTIONS", "/test", nil)
	w2 := httptest.NewRecorder()
	engine.ServeHTTP(w2, req2)

	// verify status code is 204 no content
	if w2.Code != http.StatusNoContent {
		t.Errorf("Expected status code 204 for preflight request, got %d", w2.Code)
	}
	// verify Abort() stopped the chain and core handler was not executed
	if coreHandlerExecuted {
		t.Errorf("Expected core handler not to be executed for preflight request, but it was")
	}
}
