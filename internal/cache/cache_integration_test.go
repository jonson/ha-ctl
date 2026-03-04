//go:build integration

package cache_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/jonson/ha-ctl/internal/cache"
	"github.com/jonson/ha-ctl/internal/haclient"
)

func init() {
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

func TestCacheRoundTrip(t *testing.T) {
	url := os.Getenv("HA_URL")
	token := os.Getenv("HA_TOKEN")
	if url == "" || token == "" {
		t.Fatal("HA_URL and HA_TOKEN required")
	}

	client := haclient.New(url, token)
	ctx := context.Background()

	// Refresh from HA
	c, err := cache.Refresh(ctx, client, 300)
	if err != nil {
		t.Fatalf("cache.Refresh failed: %v", err)
	}
	if len(c.Entities) == 0 {
		t.Fatal("cache has 0 entities after refresh")
	}
	t.Logf("Cached %d entities", len(c.Entities))

	// Load from disk
	loaded, err := cache.Load()
	if err != nil {
		t.Fatalf("cache.Load failed: %v", err)
	}
	if loaded == nil {
		t.Fatal("cache.Load returned nil after save")
	}
	if len(loaded.Entities) != len(c.Entities) {
		t.Errorf("entity count mismatch: saved %d, loaded %d", len(c.Entities), len(loaded.Entities))
	}
}

func TestCacheTTLExpiry(t *testing.T) {
	url := os.Getenv("HA_URL")
	token := os.Getenv("HA_TOKEN")
	if url == "" || token == "" {
		t.Fatal("HA_URL and HA_TOKEN required")
	}

	client := haclient.New(url, token)
	ctx := context.Background()

	// Refresh with TTL of 0 - should immediately be stale
	c, err := cache.Refresh(ctx, client, 0)
	if err != nil {
		t.Fatalf("cache.Refresh failed: %v", err)
	}
	if !cache.IsStale(c) {
		t.Error("cache with TTL=0 should be stale")
	}
}

func TestCacheDomainFiltering(t *testing.T) {
	url := os.Getenv("HA_URL")
	token := os.Getenv("HA_TOKEN")
	if url == "" || token == "" {
		t.Fatal("HA_URL and HA_TOKEN required")
	}

	client := haclient.New(url, token)
	ctx := context.Background()

	c, err := cache.Refresh(ctx, client, 300)
	if err != nil {
		t.Fatalf("cache.Refresh failed: %v", err)
	}

	// Filter by a domain that should exist
	lights := cache.FilterByDomain(c, "light")
	// There should be at least some entities (may or may not have lights)
	all := cache.EntityList(c)
	t.Logf("Total entities: %d, lights: %d", len(all), len(lights))

	// Verify all filtered entities are the correct domain
	for _, e := range lights {
		if e.Domain != "light" {
			t.Errorf("filtered entity %s has domain %s, expected light", e.EntityID, e.Domain)
		}
	}
}
