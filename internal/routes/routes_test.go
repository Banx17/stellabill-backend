package routes

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegister_MetricsEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	Register(router)

	// Test that /metrics endpoint exists
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/metrics", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for /metrics, got %d", w.Code)
	}

	// Verify Prometheus metrics are present
	body := w.Body.String()
	expectedMetrics := []string{
		"http_request_duration_seconds",
		"http_requests_total",
		"db_query_duration_seconds",
		"db_queries_total",
	}

	for _, metric := range expectedMetrics {
		if !strings.Contains(body, metric) {
			t.Errorf("Expected metrics output to contain %s", metric)
		}
	}
}

func TestRegister_APIEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	Register(router)

	tests := []struct {
		method   string
		path     string
		expected int
	}{
		{"GET", "/api/health", http.StatusOK},
		{"GET", "/api/subscriptions", http.StatusOK},
		{"GET", "/api/subscriptions/123", http.StatusOK},
		{"GET", "/api/plans", http.StatusOK},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(tt.method, tt.path, nil)
		router.ServeHTTP(w, req)

		if w.Code != tt.expected {
			t.Errorf("%s %s: expected status %d, got %d", tt.method, tt.path, tt.expected, w.Code)
		}
	}
}

func TestRegister_CORSHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	Register(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/api/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != 204 {
		t.Errorf("Expected status 204 for OPTIONS, got %d", w.Code)
	}
}

func TestRegister_MetricsTracksRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	Register(router)

	// Make a request to populate metrics
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	// Verify metrics were recorded
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/metrics", nil)
	router.ServeHTTP(w, req)

	body := w.Body.String()
	if !strings.Contains(body, `http_requests_total{`) {
		t.Error("Expected http_requests_total metric to be recorded")
	}
}
