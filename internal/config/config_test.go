package config

import (
	"os"
	"testing"
)

func setOrUnset(key, val string) {
	if val == "" {
		os.Unsetenv(key)
	} else {
		os.Setenv(key, val)
	}
}

func TestValidateRequiresURL(t *testing.T) {
	cfg := &Config{HAToken: "test-token", CacheTTL: 300}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing HA_URL")
	}
}

func TestValidateRequiresToken(t *testing.T) {
	cfg := &Config{HAURL: "http://localhost:8123", CacheTTL: 300}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing HA_TOKEN")
	}
}

func TestValidateRequiresHTTPScheme(t *testing.T) {
	cfg := &Config{HAURL: "ftp://localhost", HAToken: "test", CacheTTL: 300}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for non-http URL")
	}
}

func TestValidateStripsTrailingSlash(t *testing.T) {
	cfg := &Config{HAURL: "http://localhost:8123/", HAToken: "test", CacheTTL: 300}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.HAURL != "http://localhost:8123" {
		t.Errorf("expected trailing slash stripped, got %s", cfg.HAURL)
	}
}

func TestValidateDefaultsCacheTTL(t *testing.T) {
	cfg := &Config{HAURL: "http://localhost:8123", HAToken: "test", CacheTTL: -1}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.CacheTTL != DefaultCacheTTL {
		t.Errorf("expected CacheTTL=%d, got %d", DefaultCacheTTL, cfg.CacheTTL)
	}
}

func TestLoadFromEnv(t *testing.T) {
	// Save and restore original env vars
	origURL := os.Getenv("HA_URL")
	origToken := os.Getenv("HA_TOKEN")
	origTTL := os.Getenv("HA_CACHE_TTL")
	defer func() {
		setOrUnset("HA_URL", origURL)
		setOrUnset("HA_TOKEN", origToken)
		setOrUnset("HA_CACHE_TTL", origTTL)
	}()

	os.Setenv("HA_URL", "http://test:8123")
	os.Setenv("HA_TOKEN", "test-token-123")
	os.Setenv("HA_CACHE_TTL", "60")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if cfg.HAURL != "http://test:8123" {
		t.Errorf("expected URL http://test:8123, got %s", cfg.HAURL)
	}
	if cfg.HAToken != "test-token-123" {
		t.Errorf("expected token test-token-123, got %s", cfg.HAToken)
	}
	if cfg.CacheTTL != 60 {
		t.Errorf("expected TTL 60, got %d", cfg.CacheTTL)
	}
}
