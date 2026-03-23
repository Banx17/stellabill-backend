package config

import "testing"

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("ENV", "")
	t.Setenv("PORT", "")
	t.Setenv("DATABASE_URL", "")
	t.Setenv("JWT_SECRET", "")

	cfg := Load()
	if cfg.Env != "development" {
		t.Fatalf("Env: got %q want %q", cfg.Env, "development")
	}
	if cfg.Port != "8080" {
		t.Fatalf("Port: got %q want %q", cfg.Port, "8080")
	}
	if cfg.DBConn == "" {
		t.Fatalf("DBConn: expected non-empty default")
	}
	if cfg.JWTSecret != "change-me-in-production" {
		t.Fatalf("JWTSecret: got %q want %q", cfg.JWTSecret, "change-me-in-production")
	}
}

func TestLoad_Overrides(t *testing.T) {
	t.Setenv("ENV", "production")
	t.Setenv("PORT", "9999")
	t.Setenv("DATABASE_URL", "postgres://example/db")
	t.Setenv("JWT_SECRET", "secret")

	cfg := Load()
	if cfg.Env != "production" {
		t.Fatalf("Env: got %q want %q", cfg.Env, "production")
	}
	if cfg.Port != "9999" {
		t.Fatalf("Port: got %q want %q", cfg.Port, "9999")
	}
	if cfg.DBConn != "postgres://example/db" {
		t.Fatalf("DBConn: got %q want %q", cfg.DBConn, "postgres://example/db")
	}
	if cfg.JWTSecret != "secret" {
		t.Fatalf("JWTSecret: got %q want %q", cfg.JWTSecret, "secret")
	}
}
