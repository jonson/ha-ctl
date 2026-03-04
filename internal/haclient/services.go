package haclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// CallService calls a Home Assistant service.
func (c *Client) CallService(ctx context.Context, domain, service string, data map[string]interface{}) ([]Entity, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshaling service data: %w", err)
	}

	respData, err := c.doRequest(ctx, "POST", fmt.Sprintf("/api/services/%s/%s", domain, service), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("calling service %s.%s: %w", domain, service, err)
	}

	var entities []Entity
	if err := json.Unmarshal(respData, &entities); err != nil {
		return nil, fmt.Errorf("parsing service response: %w", err)
	}
	return entities, nil
}
