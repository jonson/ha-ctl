package config

import (
	"fmt"
	"strings"
)

// Validate checks that required config fields are set and valid.
func (c *Config) Validate() error {
	if c.HAURL == "" {
		return fmt.Errorf("HA_URL is required (set via environment variable or config file)")
	}
	if !strings.HasPrefix(c.HAURL, "http://") && !strings.HasPrefix(c.HAURL, "https://") {
		return fmt.Errorf("HA_URL must start with http:// or https://")
	}
	// Strip trailing slash for consistency
	c.HAURL = strings.TrimRight(c.HAURL, "/")

	if c.HAToken == "" {
		return fmt.Errorf("HA_TOKEN is required (set via environment variable or config file)")
	}
	if c.CacheTTL <= 0 {
		c.CacheTTL = DefaultCacheTTL
	}
	return nil
}
