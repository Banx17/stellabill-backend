package config

import (
	"os"
	"testing"
)

func TestLoad_DefaultValues(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("ENV")
	os.Unsetenv("PORT")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("JWT_SECRET")

	cfg := Load()

	if cfg.Env != "development" {
		t.Errorf("Expected Env to be 'development', got %s", cfg.Env)
	}

	if cfg.Port != "8080" {
		t.Errorf("Expected Port to be '8080', got %s", cfg.Port)
	}

	expectedDB := "postgres://localhost/stellarbill?sslmode=disable"
	if cfg.DBConn != expectedDB {
		t.Errorf("Expected DBConn to be '%s', got %s", expectedDB, cfg.DBConn)
	}

	if cfg.JWTSecret != "change-me-in-production" {
		t.Errorf("Expected JWTSecret to be 'change-me-in-production', got %s", cfg.JWTSecret)
	}
}

func TestLoad_CustomValues(t *testing.T) {
	// Set custom environment variables
	os.Setenv("ENV", "production")
	os.Setenv("PORT", "3000")
	os.Setenv("DATABASE_URL", "postgres://custom/db")
	os.Setenv("JWT_SECRET", "my-secret")

	// Clear after test
	defer func() {
		os.Unsetenv("ENV")
		os.Unsetenv("PORT")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("JWT_SECRET")
	}()

	cfg := Load()

	if cfg.Env != "production" {
		t.Errorf("Expected Env to be 'production', got %s", cfg.Env)
	}

	if cfg.Port != "3000" {
		t.Errorf("Expected Port to be '3000', got %s", cfg.Port)
	}

	if cfg.DBConn != "postgres://custom/db" {
		t.Errorf("Expected DBConn to be 'postgres://custom/db', got %s", cfg.DBConn)
	}

	if cfg.JWTSecret != "my-secret" {
		t.Errorf("Expected JWTSecret to be 'my-secret', got %s", cfg.JWTSecret)
	}
}

func TestLoad_PORT_Override(t *testing.T) {
	// Set PORT via environment
	os.Setenv("PORT", "9090")
	defer os.Unsetenv("PORT")

	cfg := Load()

	if cfg.Port != "9090" {
		t.Errorf("Expected Port to be '9090', got %s", cfg.Port)
	}
}
