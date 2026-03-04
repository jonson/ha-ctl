package config

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v3"

	"github.com/jonson/ha-ctl/internal/util"
)

// Config holds all ha-ctl configuration.
type Config struct {
	HAURL    string `yaml:"ha_url"`
	HAToken  string `yaml:"ha_token"`
	CacheTTL int    `yaml:"cache_ttl"`
}

// Load reads configuration with precedence: env vars > config file > defaults.
func Load() (*Config, error) {
	cfg := &Config{
		CacheTTL: DefaultCacheTTL,
	}

	// Load from config file first (lowest precedence after defaults)
	loadFromFile(cfg)

	// Environment variables override config file
	if v := os.Getenv("HA_URL"); v != "" {
		cfg.HAURL = v
	}
	if v := os.Getenv("HA_TOKEN"); v != "" {
		cfg.HAToken = v
	}
	if v := os.Getenv("HA_CACHE_TTL"); v != "" {
		if ttl, err := strconv.Atoi(v); err == nil {
			cfg.CacheTTL = ttl
		}
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// loadFromFile reads the YAML config file if it exists.
func loadFromFile(cfg *Config) {
	path := util.ConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return // file doesn't exist, that's fine
	}
	_ = yaml.Unmarshal(data, cfg) // intentional: partial/malformed config is OK, env vars override
}
