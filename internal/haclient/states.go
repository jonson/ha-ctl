package haclient

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetStates fetches all entity states from Home Assistant.
func (c *Client) GetStates(ctx context.Context) ([]Entity, error) {
	data, err := c.doRequest(ctx, "GET", "/api/states", nil)
	if err != nil {
		return nil, fmt.Errorf("fetching states: %w", err)
	}

	var entities []Entity
	if err := json.Unmarshal(data, &entities); err != nil {
		return nil, fmt.Errorf("parsing states response: %w", err)
	}
	return entities, nil
}

// GetState fetches the state of a single entity.
func (c *Client) GetState(ctx context.Context, entityID string) (*Entity, error) {
	data, err := c.doRequest(ctx, "GET", "/api/states/"+entityID, nil)
	if err != nil {
		return nil, fmt.Errorf("fetching state for %s: %w", entityID, err)
	}

	var entity Entity
	if err := json.Unmarshal(data, &entity); err != nil {
		return nil, fmt.Errorf("parsing state response: %w", err)
	}
	return &entity, nil
}
