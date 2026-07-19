package goroute

import (
	"bytes"
	"log"
	"net/http/httptest"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	// hijack the standard logger to capture the output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr) //restore the original output after the test

	engine := New()
	engine.Use(Logger()) // Apply the Logger middleware

	engine.GET("/ping", func(c *Context) {
		c.String(200, "pong")
	})

	req := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("[GoExpress]")) {
		t.Errorf("Expected log output to contain '[GoExpress]', got: %s", output)
	}

}
