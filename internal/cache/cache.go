package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jonson/ha-ctl/internal/util"
)

// Load reads the cache from disk. Returns nil if no cache exists.
func Load() (*Cache, error) {
	path := util.CachePath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading cache: %w", err)
	}

	var c Cache
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parsing cache: %w", err)
	}
	return &c, nil
}

// Save writes the cache to disk atomically (temp file + rename).
func Save(c *Cache) error {
	dir := util.CacheDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating cache directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling cache: %w", err)
	}

	path := util.CachePath()
	tmpPath := path + ".tmp"

	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return fmt.Errorf("writing temp cache file: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("renaming cache file: %w", err)
	}
	return nil
}

// IsStale returns true if the cache is older than its TTL.
func IsStale(c *Cache) bool {
	if c == nil {
		return true
	}
	age := time.Since(c.Timestamp)
	return age.Seconds() > float64(c.TTLSeconds)
}

// FilterByDomain returns only entities matching the given domain.
func FilterByDomain(c *Cache, domain string) []CacheEntity {
	var result []CacheEntity
	for _, e := range c.Entities {
		if e.Domain == domain {
			result = append(result, e)
		}
	}
	return result
}

// EntityList returns all cache entities as a slice, sorted doesn't matter.
func EntityList(c *Cache) []CacheEntity {
	result := make([]CacheEntity, 0, len(c.Entities))
	for _, e := range c.Entities {
		result = append(result, e)
	}
	return result
}

// Search returns entities matching query (case-insensitive substring) against
// entity_id and friendly_name. If domain is non-empty, results are also
// filtered to that domain.
func Search(c *Cache, query string, domain string) []CacheEntity {
	q := strings.ToLower(query)
	var result []CacheEntity
	for _, e := range c.Entities {
		if domain != "" && e.Domain != domain {
			continue
		}
		if strings.Contains(strings.ToLower(e.EntityID), q) ||
			strings.Contains(strings.ToLower(e.FriendlyName), q) {
			result = append(result, e)
		}
	}
	return result
}

// FilterByState returns only entities whose state matches (case-insensitive).
func FilterByState(entities []CacheEntity, state string) []CacheEntity {
	s := strings.ToLower(state)
	var result []CacheEntity
	for _, e := range entities {
		if strings.ToLower(e.State) == s {
			result = append(result, e)
		}
	}
	return result
}

// Dir returns the cache directory path (exposed for commands that need it).
func Dir() string {
	return util.CacheDir()
}

// Path returns the cache file path (exposed for commands that need it).
func Path() string {
	return util.CachePath()
}

// AbsProjectRoot attempts to find the project root for .env loading in tests.
func AbsProjectRoot() string {
	// Walk up from current directory looking for go.mod
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
