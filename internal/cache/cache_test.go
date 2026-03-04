package cache

import (
	"testing"
	"time"
)

func TestIsStaleNilCache(t *testing.T) {
	if !IsStale(nil) {
		t.Error("nil cache should be stale")
	}
}

func TestIsStaleExpired(t *testing.T) {
	c := &Cache{
		Timestamp:  time.Now().Add(-10 * time.Minute),
		TTLSeconds: 300, // 5 minutes
	}
	if !IsStale(c) {
		t.Error("10-minute-old cache with 5-minute TTL should be stale")
	}
}

func TestIsStaleStillFresh(t *testing.T) {
	c := &Cache{
		Timestamp:  time.Now(),
		TTLSeconds: 300,
	}
	if IsStale(c) {
		t.Error("just-created cache should not be stale")
	}
}

func TestIsStaleZeroTTL(t *testing.T) {
	c := &Cache{
		Timestamp:  time.Now(),
		TTLSeconds: 0,
	}
	// With TTL 0, any non-zero age makes it stale
	// However time.Since(time.Now()) might be 0, so this is an edge case.
	// We just test it doesn't panic.
	_ = IsStale(c)
}

func TestFilterByDomain(t *testing.T) {
	c := &Cache{
		Entities: map[string]CacheEntity{
			"light.kitchen":   {EntityID: "light.kitchen", Domain: "light"},
			"light.bedroom":   {EntityID: "light.bedroom", Domain: "light"},
			"switch.garage":   {EntityID: "switch.garage", Domain: "switch"},
			"sensor.temp":     {EntityID: "sensor.temp", Domain: "sensor"},
		},
	}

	lights := FilterByDomain(c, "light")
	if len(lights) != 2 {
		t.Errorf("expected 2 lights, got %d", len(lights))
	}

	switches := FilterByDomain(c, "switch")
	if len(switches) != 1 {
		t.Errorf("expected 1 switch, got %d", len(switches))
	}

	empty := FilterByDomain(c, "climate")
	if len(empty) != 0 {
		t.Errorf("expected 0 climate entities, got %d", len(empty))
	}
}

func TestEntityList(t *testing.T) {
	c := &Cache{
		Entities: map[string]CacheEntity{
			"light.a": {EntityID: "light.a"},
			"light.b": {EntityID: "light.b"},
		},
	}
	list := EntityList(c)
	if len(list) != 2 {
		t.Errorf("expected 2 entities, got %d", len(list))
	}
}

func TestExtractKeyAttributes(t *testing.T) {
	attrs := map[string]any{
		"brightness":   128,
		"color_temp":   350,
		"friendly_name": "Kitchen Light",
		"supported_features": 44,
	}

	result := extractKeyAttributes("light", attrs)
	if result == nil {
		t.Fatal("expected non-nil key attributes for light")
	}
	if result["brightness"] != 128 {
		t.Errorf("expected brightness=128, got %v", result["brightness"])
	}
	if result["color_temp"] != 350 {
		t.Errorf("expected color_temp=350, got %v", result["color_temp"])
	}
	if _, exists := result["friendly_name"]; exists {
		t.Error("friendly_name should not be in key attributes")
	}
	if _, exists := result["supported_features"]; exists {
		t.Error("supported_features should not be in key attributes")
	}
}

func TestExtractKeyAttributesUnknownDomain(t *testing.T) {
	attrs := map[string]any{"foo": "bar"}
	result := extractKeyAttributes("unknown_domain", attrs)
	if result != nil {
		t.Errorf("expected nil for unknown domain, got %v", result)
	}
}

func TestSearch(t *testing.T) {
	c := &Cache{
		Entities: map[string]CacheEntity{
			"light.kitchen_main":  {EntityID: "light.kitchen_main", FriendlyName: "Kitchen Main Light", Domain: "light"},
			"light.bedroom":      {EntityID: "light.bedroom", FriendlyName: "Bedroom Light", Domain: "light"},
			"switch.kitchen_fan": {EntityID: "switch.kitchen_fan", FriendlyName: "Kitchen Fan", Domain: "switch"},
			"sensor.temp":        {EntityID: "sensor.temp", FriendlyName: "Temperature", Domain: "sensor"},
		},
	}

	// Search by name, no domain filter
	results := Search(c, "kitchen", "")
	if len(results) != 2 {
		t.Errorf("expected 2 kitchen matches, got %d", len(results))
	}

	// Search with domain filter
	results = Search(c, "kitchen", "light")
	if len(results) != 1 {
		t.Errorf("expected 1 kitchen light, got %d", len(results))
	}
	if len(results) > 0 && results[0].EntityID != "light.kitchen_main" {
		t.Errorf("expected light.kitchen_main, got %s", results[0].EntityID)
	}

	// Case-insensitive search
	results = Search(c, "BEDROOM", "")
	if len(results) != 1 {
		t.Errorf("expected 1 bedroom match, got %d", len(results))
	}

	// Search matching entity_id but not friendly_name
	results = Search(c, "temp", "")
	if len(results) != 1 {
		t.Errorf("expected 1 temp match, got %d", len(results))
	}

	// No matches
	results = Search(c, "nonexistent", "")
	if len(results) != 0 {
		t.Errorf("expected 0 matches, got %d", len(results))
	}
}

func TestFilterByState(t *testing.T) {
	entities := []CacheEntity{
		{EntityID: "light.a", State: "on"},
		{EntityID: "light.b", State: "off"},
		{EntityID: "light.c", State: "on"},
		{EntityID: "light.d", State: "unavailable"},
	}

	on := FilterByState(entities, "on")
	if len(on) != 2 {
		t.Errorf("expected 2 'on' entities, got %d", len(on))
	}

	off := FilterByState(entities, "off")
	if len(off) != 1 {
		t.Errorf("expected 1 'off' entity, got %d", len(off))
	}

	// Case-insensitive
	unavail := FilterByState(entities, "UNAVAILABLE")
	if len(unavail) != 1 {
		t.Errorf("expected 1 unavailable entity, got %d", len(unavail))
	}

	// No matches
	none := FilterByState(entities, "idle")
	if len(none) != 0 {
		t.Errorf("expected 0 'idle' entities, got %d", len(none))
	}
}
