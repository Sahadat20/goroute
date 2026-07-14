package goroute

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Define a dummy struct for testing the reflection engine
type TestUser struct {
	Username string `json:"username" validate:"required"`
	Age      int    `json:"age" validate:"required"`
	Bio      string `json:"bio"` // Not required
}

func TestDataBindingAndValidation(t *testing.T) {
	engine := New()

	// Register a test route that attempts to bind incoming JSON
	engine.POST("/register", func(c *Context) {
		var user TestUser

		// Attempt to bind and catch validation errors
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, map[string]string{"status": "success", "user": user.Username})
	})

	// --- SUBTEST 1: Valid Payload ---
	validJSON := []byte(`{"username":"system_admin", "age":35, "bio":"Hello world"}`)
	req1 := httptest.NewRequest("POST", "/register", bytes.NewBuffer(validJSON))
	w1 := httptest.NewRecorder()
	engine.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("Subtest 1 Failed: Expected 200 OK, got %d", w1.Code)
	}

	// --- SUBTEST 2: Invalid JSON (Syntax Error) ---
	badJSON := []byte(`{"username":"system_admin", "age":}`) // Broken JSON syntax
	req2 := httptest.NewRequest("POST", "/register", bytes.NewBuffer(badJSON))
	w2 := httptest.NewRecorder()
	engine.ServeHTTP(w2, req2)

	if w2.Code != http.StatusBadRequest {
		t.Errorf("Subtest 2 Failed: Expected 400 Bad Request for malformed JSON, got %d", w2.Code)
	}

	// --- SUBTEST 3: Validation Failure (Missing Required Field) ---
	missingFieldJSON := []byte(`{"username":"system_admin"}`) // Age is missing
	req3 := httptest.NewRequest("POST", "/register", bytes.NewBuffer(missingFieldJSON))
	w3 := httptest.NewRecorder()
	engine.ServeHTTP(w3, req3)

	if w3.Code != http.StatusBadRequest {
		t.Errorf("Subtest 3 Failed: Expected 400 Bad Request for missing required field, got %d", w3.Code)
	}
}
