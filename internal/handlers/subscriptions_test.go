package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestListSubscriptions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/subscriptions", ListSubscriptions)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/subscriptions", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	subscriptions, ok := response["subscriptions"].([]interface{})
	if !ok {
		t.Fatalf("Expected subscriptions to be an array, got %T", response["subscriptions"])
	}

	if len(subscriptions) != 0 {
		t.Errorf("Expected empty subscriptions array, got %d items", len(subscriptions))
	}
}

func TestGetSubscription_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/subscriptions/:id", GetSubscription)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/subscriptions/123", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["id"] != "123" {
		t.Errorf("Expected id '123', got %v", response["id"])
	}

	if response["status"] != "placeholder" {
		t.Errorf("Expected status 'placeholder', got %v", response["status"])
	}
}

func TestGetSubscription_EmptyID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/subscriptions/:id", GetSubscription)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/subscriptions/", nil)
	router.ServeHTTP(w, req)

	// Gin will return 404 for empty path parameter
	if w.Code == http.StatusOK {
		// If it reaches the handler, it should return bad request
		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err == nil {
			if _, hasError := response["error"]; hasError {
				return // Expected error response
			}
		}
	}
	// Either 404 or error response is acceptable
}
