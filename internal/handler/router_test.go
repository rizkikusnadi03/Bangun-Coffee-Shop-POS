package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewRouterHealthEndpoint(t *testing.T) {
	router := NewRouter("development")
	if gin.Mode() != gin.DebugMode {
		t.Errorf("expected debug mode, got: %s", gin.Mode())
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/v1/health", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse json response: %v", err)
	}

	if resp["success"] != true {
		t.Errorf("expected success: true, got %v", resp["success"])
	}

	if resp["message"] != "server is running" {
		t.Errorf("expected message: 'server is running', got %v", resp["message"])
	}
}

func TestNewRouterProductionMode(t *testing.T) {
	_ = NewRouter("production")
	if gin.Mode() != gin.ReleaseMode {
		t.Errorf("expected release mode, got: %s", gin.Mode())
	}
}
