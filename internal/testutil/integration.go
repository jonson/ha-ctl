//go:build integration

package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/jonson/ha-ctl/internal/haclient"
)

// LoadEnv loads .env from the project root and returns HA credentials.
// Fails the test immediately if any required variable is missing.
func LoadEnv(t *testing.T) (url, token, testEntity string) {
	t.Helper()
	// Try to find and load .env from project root
	root := findProjectRoot()
	if root != "" {
		_ = godotenv.Load(filepath.Join(root, ".env"))
	}
	url = os.Getenv("HA_URL")
	token = os.Getenv("HA_TOKEN")
	testEntity = os.Getenv("HA_TEST_ENTITY")
	if url == "" || token == "" || testEntity == "" {
		t.Fatal("integration test requires HA_URL, HA_TOKEN, HA_TEST_ENTITY in .env or environment")
	}
	return
}

// NewTestClient returns a configured haclient.Client for integration tests.
func NewTestClient(t *testing.T) (*haclient.Client, string) {
	t.Helper()
	url, token, testEntity := LoadEnv(t)
	client := haclient.New(url, token)
	return client, testEntity
}

// findProjectRoot walks up directories looking for go.mod.
func findProjectRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}
