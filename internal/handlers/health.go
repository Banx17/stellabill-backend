package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"stellarbill-backend/internal/outbox"
)

var (
	globalOutboxManager *outbox.Manager
)

// SetOutboxManager sets the global outbox manager for health checks
func SetOutboxManager(manager *outbox.Manager) {
	globalOutboxManager = manager
}

func Health(c *gin.Context) {
	status := gin.H{
		"status":  "ok",
		"service": "stellarbill-backend",
	}

	// Check outbox health if available
	if globalOutboxManager != nil {
		if err := globalOutboxManager.Health(); err != nil {
			status["status"] = "degraded"
			status["outbox"] = gin.H{
				"status": "unhealthy",
				"error":  err.Error(),
			}
		} else {
			stats, err := globalOutboxManager.GetStats()
			if err == nil {
				status["outbox"] = stats
			}
		}
	}

	c.JSON(http.StatusOK, status)
}

// OutboxStats returns detailed outbox statistics
func OutboxStats(c *gin.Context) {
	if globalOutboxManager == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Outbox manager not available",
		})
		return
	}

	stats, err := globalOutboxManager.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// PublishTestEvent publishes a test event for development/testing
func PublishTestEvent(c *gin.Context) {
	if globalOutboxManager == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Outbox manager not available",
		})
		return
	}

	// Get event type from query parameter
	eventType := c.Query("type")
	if eventType == "" {
		eventType = "test.event"
	}

	// Create test event data
	eventData := gin.H{
		"message":     "This is a test event",
		"timestamp":   gin.H{"$date": gin.H{"$numberLong": strconv.FormatInt(c.Request.Context().Value("timestamp").(int64), 10)}},
		"request_id":  c.GetHeader("X-Request-ID"),
		"user_agent":  c.GetHeader("User-Agent"),
		"ip_address": c.ClientIP(),
	}

	// Publish the event
	service := globalOutboxManager.GetService()
	err := service.PublishEvent(c.Request.Context(), eventType, eventData, nil, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Test event published successfully",
		"event_type": eventType,
	})
}
