package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env       string
	Port      string
	DBConn    string
	JWTSecret string
	Outbox    OutboxConfig
}

type OutboxConfig struct {
	PublisherType       string
	HTTPEndpoint        string
	PollInterval        string
	BatchSize           int
	MaxRetries          int
	RetryBackoffFactor  float64
	CleanupInterval     string
	CompletedEventTTL   string
	ProcessingTimeout   string
}

func Load() Config {
	return Config{
		Env:       getEnv("ENV", "development"),
		Port:      getEnv("PORT", "8080"),
		DBConn:    getEnv("DATABASE_URL", "postgres://localhost/stellarbill?sslmode=disable"),
		JWTSecret: getEnv("JWT_SECRET", "change-me-in-production"),
		Outbox: OutboxConfig{
			PublisherType:       getEnv("OUTBOX_PUBLISHER_TYPE", "console"),
			HTTPEndpoint:        getEnv("OUTBOX_HTTP_ENDPOINT", ""),
			PollInterval:        getEnv("OUTBOX_POLL_INTERVAL", "5s"),
			BatchSize:           getIntEnv("OUTBOX_BATCH_SIZE", 10),
			MaxRetries:          getIntEnv("OUTBOX_MAX_RETRIES", 3),
			RetryBackoffFactor:  getFloatEnv("OUTBOX_RETRY_BACKOFF_FACTOR", 2.0),
			CleanupInterval:     getEnv("OUTBOX_CLEANUP_INTERVAL", "1h"),
			CompletedEventTTL:   getEnv("OUTBOX_COMPLETED_EVENT_TTL", "24h"),
			ProcessingTimeout:   getEnv("OUTBOX_PROCESSING_TIMEOUT", "30s"),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getIntEnv(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getFloatEnv(key string, fallback float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return fallback
}

func (c OutboxConfig) GetPollInterval() time.Duration {
	if d, err := time.ParseDuration(c.PollInterval); err == nil {
		return d
	}
	return 5 * time.Second
}

func (c OutboxConfig) GetCleanupInterval() time.Duration {
	if d, err := time.ParseDuration(c.CleanupInterval); err == nil {
		return d
	}
	return 1 * time.Hour
}

func (c OutboxConfig) GetCompletedEventTTL() time.Duration {
	if d, err := time.ParseDuration(c.CompletedEventTTL); err == nil {
		return d
	}
	return 24 * time.Hour
}

func (c OutboxConfig) GetProcessingTimeout() time.Duration {
	if d, err := time.ParseDuration(c.ProcessingTimeout); err == nil {
		return d
	}
	return 30 * time.Second
}
