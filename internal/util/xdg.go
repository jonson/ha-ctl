package util

import (
	"os"
	"path/filepath"
)

// ConfigDir returns the ha-ctl config directory (~/.config/ha-ctl).
func ConfigDir() string {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, _ := os.UserHomeDir()
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "ha-ctl")
}

// CacheDir returns the ha-ctl cache directory (~/.cache/ha-ctl).
func CacheDir() string {
	base := os.Getenv("XDG_CACHE_HOME")
	if base == "" {
		home, _ := os.UserHomeDir()
		base = filepath.Join(home, ".cache")
	}
	return filepath.Join(base, "ha-ctl")
}

// CachePath returns the full path to the entities cache file.
func CachePath() string {
	return filepath.Join(CacheDir(), "entities.json")
}

// ConfigPath returns the full path to the config file.
func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config.yaml")
}
