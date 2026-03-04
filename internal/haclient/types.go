package haclient

// Entity represents a Home Assistant entity state.
type Entity struct {
	EntityID    string                 `json:"entity_id"`
	State       string                 `json:"state"`
	Attributes  map[string]interface{} `json:"attributes"`
	LastChanged string                 `json:"last_changed"`
	LastUpdated string                 `json:"last_updated"`
}

// FriendlyName returns the friendly_name attribute or the entity_id as fallback.
func (e *Entity) FriendlyName() string {
	if name, ok := e.Attributes["friendly_name"].(string); ok {
		return name
	}
	return e.EntityID
}

// Domain extracts the domain from the entity_id (e.g. "light" from "light.kitchen").
func (e *Entity) Domain() string {
	for i, c := range e.EntityID {
		if c == '.' {
			return e.EntityID[:i]
		}
	}
	return e.EntityID
}

// ServiceCall represents a request to call a Home Assistant service.
type ServiceCall struct {
	Domain  string                 `json:"domain"`
	Service string                 `json:"service"`
	Data    map[string]interface{} `json:"service_data,omitempty"`
}
