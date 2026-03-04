package cache

import "time"

// Cache represents the cached entity data.
type Cache struct {
	Version    string                 `json:"version"`
	Timestamp  time.Time              `json:"timestamp"`
	TTLSeconds int                    `json:"ttl_seconds"`
	Entities   map[string]CacheEntity `json:"entities"`
}

// CacheEntity represents a cached entity with key attributes.
type CacheEntity struct {
	EntityID      string                 `json:"entity_id"`
	FriendlyName  string                 `json:"friendly_name"`
	Domain        string                 `json:"domain"`
	State         string                 `json:"state"`
	KeyAttributes map[string]interface{} `json:"key_attributes"`
}

// keyAttributesByDomain defines which attributes to cache per domain.
var keyAttributesByDomain = map[string][]string{
	"light":        {"brightness", "color_temp", "rgb_color", "effect"},
	"climate":      {"temperature", "current_temperature", "hvac_mode"},
	"media_player": {"volume_level", "media_title", "source"},
	"cover":        {"current_position"},
	"sensor":       {"unit_of_measurement", "device_class"},
	"switch":       {},
}
