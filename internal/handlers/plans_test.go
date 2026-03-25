package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestListPlans(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/plans", ListPlans)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/plans", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	plans, ok := response["plans"].([]interface{})
	if !ok {
		t.Fatalf("Expected plans to be an array, got %T", response["plans"])
	}

	if len(plans) != 0 {
		t.Errorf("Expected empty plans array, got %d items", len(plans))
	}
}
