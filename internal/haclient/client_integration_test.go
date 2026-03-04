//go:build integration

package haclient_test

import (
	"context"
	"testing"

	"github.com/jonson/ha-ctl/internal/haclient"
	"github.com/jonson/ha-ctl/internal/testutil"
)

func TestAPIConnectivity(t *testing.T) {
	client, _ := testutil.NewTestClient(t)
	ctx := context.Background()

	// GET /api/states should work
	states, err := client.GetStates(ctx)
	if err != nil {
		t.Fatalf("GetStates failed: %v", err)
	}
	if len(states) == 0 {
		t.Fatal("GetStates returned 0 entities")
	}
	t.Logf("GetStates returned %d entities", len(states))
}

func TestInvalidToken(t *testing.T) {
	url, _, _ := testutil.LoadEnv(t)
	client := haclient.New(url, "invalid-token-12345")
	ctx := context.Background()

	_, err := client.GetStates(ctx)
	if err == nil {
		t.Fatal("expected error with invalid token, got nil")
	}
	t.Logf("Got expected error: %v", err)
}

func TestGetState(t *testing.T) {
	client, testEntity := testutil.NewTestClient(t)
	ctx := context.Background()

	entity, err := client.GetState(ctx, testEntity)
	if err != nil {
		t.Fatalf("GetState(%s) failed: %v", testEntity, err)
	}
	if entity.EntityID != testEntity {
		t.Errorf("expected entity_id %s, got %s", testEntity, entity.EntityID)
	}
	if entity.State == "" {
		t.Error("entity state is empty")
	}
	t.Logf("Entity %s state: %s", entity.EntityID, entity.State)
}

func TestGetStateNotFound(t *testing.T) {
	client, _ := testutil.NewTestClient(t)
	ctx := context.Background()

	_, err := client.GetState(ctx, "nonexistent.entity_that_does_not_exist")
	if err == nil {
		t.Fatal("expected error for nonexistent entity, got nil")
	}
	t.Logf("Got expected error: %v", err)
}
