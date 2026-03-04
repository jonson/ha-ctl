package cache

import (
	"context"
	"time"

	"github.com/jonson/ha-ctl/internal/haclient"
)

// Refresh fetches all states from the HA API and rebuilds the cache.
func Refresh(ctx context.Context, client *haclient.Client, ttlSeconds int) (*Cache, error) {
	entities, err := client.GetStates(ctx)
	if err != nil {
		return nil, err
	}

	c := &Cache{
		Version:    "v1",
		Timestamp:  time.Now(),
		TTLSeconds: ttlSeconds,
		Entities:   make(map[string]CacheEntity, len(entities)),
	}

	for _, e := range entities {
		domain := e.Domain()
		ce := CacheEntity{
			EntityID:      e.EntityID,
			FriendlyName:  e.FriendlyName(),
			Domain:        domain,
			State:         e.State,
			KeyAttributes: extractKeyAttributes(domain, e.Attributes),
		}
		c.Entities[e.EntityID] = ce
	}

	if err := Save(c); err != nil {
		return nil, err
	}

	return c, nil
}

// extractKeyAttributes pulls only the relevant attributes for a given domain.
func extractKeyAttributes(domain string, attrs map[string]interface{}) map[string]interface{} {
	keys, ok := keyAttributesByDomain[domain]
	if !ok || len(keys) == 0 {
		return nil
	}

	result := make(map[string]interface{})
	for _, key := range keys {
		if val, exists := attrs[key]; exists {
			result[key] = val
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}
