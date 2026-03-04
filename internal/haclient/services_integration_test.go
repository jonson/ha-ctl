//go:build integration

package haclient_test

import (
	"context"
	"testing"
	"time"

	"github.com/jonson/ha-ctl/internal/testutil"
)

func TestCallServiceToggle(t *testing.T) {
	client, testEntity := testutil.NewTestClient(t)
	ctx := context.Background()

	// Get initial state
	initial, err := client.GetState(ctx, testEntity)
	if err != nil {
		t.Fatalf("GetState failed: %v", err)
	}
	t.Logf("Initial state: %s", initial.State)

	// Turn on
	_, err = client.CallService(ctx, "input_boolean", "turn_on", map[string]any{
		"entity_id": testEntity,
	})
	if err != nil {
		t.Fatalf("CallService turn_on failed: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	entity, err := client.GetState(ctx, testEntity)
	if err != nil {
		t.Fatalf("GetState after turn_on failed: %v", err)
	}
	if entity.State != "on" {
		t.Errorf("expected state 'on' after turn_on, got '%s'", entity.State)
	}

	// Turn off
	_, err = client.CallService(ctx, "input_boolean", "turn_off", map[string]any{
		"entity_id": testEntity,
	})
	if err != nil {
		t.Fatalf("CallService turn_off failed: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	entity, err = client.GetState(ctx, testEntity)
	if err != nil {
		t.Fatalf("GetState after turn_off failed: %v", err)
	}
	if entity.State != "off" {
		t.Errorf("expected state 'off' after turn_off, got '%s'", entity.State)
	}

	// Restore original state
	restoreService := "turn_off"
	if initial.State == "on" {
		restoreService = "turn_on"
	}
	_, _ = client.CallService(ctx, "input_boolean", restoreService, map[string]any{
		"entity_id": testEntity,
	})
}
