//go:build integration

package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/jonson/ha-ctl/internal/config"
)

func init() {
	// Find and load .env from project root
	dir, _ := os.Getwd()
	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			_ = godotenv.Load(envPath)
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
}

func TestConfigLoadsFromEnv(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load() failed: %v", err)
	}
	if cfg.HAURL == "" {
		t.Error("HAURL is empty")
	}
	if cfg.HAToken == "" {
		t.Error("HAToken is empty")
	}
	if cfg.CacheTTL <= 0 {
		t.Error("CacheTTL should be positive")
	}
	t.Logf("Config loaded: URL=%s, TTL=%d", cfg.HAURL, cfg.CacheTTL)
}

func TestEnvVarPrecedence(t *testing.T) {
	original := os.Getenv("HA_CACHE_TTL")
	defer os.Setenv("HA_CACHE_TTL", original)

	os.Setenv("HA_CACHE_TTL", "999")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load() failed: %v", err)
	}
	if cfg.CacheTTL != 999 {
		t.Errorf("expected CacheTTL=999 from env override, got %d", cfg.CacheTTL)
	}
}
