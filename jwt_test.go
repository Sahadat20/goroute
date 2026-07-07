package goroute

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TESTJWTFrameWorkIntegration(t *testing.T) {
	secret := "this-is-my-secure-key"
	engine := New()

	adminGroup := engine.Group("/")
	adminGroup.Use(JWTAuth(secret))

	adminGroup.GET("/dashboard", func(c *Context) {
		user, _ := c.Get("username")
		c.String(http.StatusOK, "Welcome "+user.(string))
	})

	// SUBTEST 1: Request lacking token
	req1 := httptest.NewRequest("GET", "/dashboard", nil)
	w1 := httptest.NewRecorder()
	engine.ServeHTTP(w1, req1)
	if w1.Code != http.StatusUnauthorized {
		t.Errorf("Subtest 1 Failed: Expected status 401, got %d", w1.Code)
	}
	token, _ := GenerateJWT("dev_user", secret, 5*time.Minute)
	req2 := httptest.NewRequest("GET", "/dashboard", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	engine.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Errorf("Subtest 2 Failed: Expected status 200, got %d", w2.Code)
	}
	if w2.Body.String() != "Welcome dev_user" {
		t.Errorf("Expected greeting 'Welcome dev_user', got '%s'", w2.Body.String())
	}

	// --- SUBTEST 3: Request with Tampered Token Signature ---
	tamperedToken := token + "manipulation123"
	req3 := httptest.NewRequest("GET", "/dashboard", nil)
	req3.Header.Set("Authorization", "Bearer "+tamperedToken)
	w3 := httptest.NewRecorder()
	engine.ServeHTTP(w3, req3)

	if w3.Code != http.StatusUnauthorized {
		t.Errorf("Subtest 3 Failed: Core handler executed instead of rejecting signature anomaly!")
	}
}
