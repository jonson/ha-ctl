package output

import (
	"fmt"
	"strings"
)

// TextFormatter outputs data as human-readable text.
type TextFormatter struct{}

// Format formats data as text. It handles known types with custom formatting
// and falls back to fmt.Sprintf for unknown types.
func (f *TextFormatter) Format(data interface{}) (string, error) {
	switch v := data.(type) {
	case EntityListOutput:
		return formatEntityList(v), nil
	case EntityStateOutput:
		return formatEntityState(v), nil
	case ServiceCallOutput:
		return formatServiceCall(v), nil
	case ContextOutput:
		return formatContext(v), nil
	case CacheStatsOutput:
		return formatCacheStats(v), nil
	default:
		return fmt.Sprintf("%v", data), nil
	}
}

// EntityListOutput is the output type for the entities command.
type EntityListOutput struct {
	Entities []EntityItem `json:"entities"`
	Count    int          `json:"count"`
}

// EntityItem represents a single entity in list output.
type EntityItem struct {
	EntityID      string                 `json:"entity_id"`
	State         string                 `json:"state"`
	FriendlyName  string                 `json:"friendly_name"`
	Domain        string                 `json:"domain"`
	KeyAttributes map[string]interface{} `json:"key_attributes,omitempty"`
}

// EntityStateOutput is the output type for the state command.
type EntityStateOutput struct {
	EntityID    string                 `json:"entity_id"`
	State       string                 `json:"state"`
	Attributes  map[string]interface{} `json:"attributes"`
	LastChanged string                 `json:"last_changed"`
	LastUpdated string                 `json:"last_updated"`
}

// ServiceCallOutput is the output type for the call command.
type ServiceCallOutput struct {
	Success  bool   `json:"success"`
	Domain   string `json:"domain"`
	Service  string `json:"service"`
	EntityID string `json:"entity_id,omitempty"`
}

// ContextOutput is the output type for the context command.
type ContextOutput struct {
	Summary      string          `json:"summary"`
	Controllable []DomainSummary `json:"controllable"`
	Other        map[string]int  `json:"other,omitempty"`
}

// DomainSummary is a summary of entities in a domain.
type DomainSummary struct {
	Domain   string              `json:"domain"`
	Count    int                 `json:"count"`
	Entities []EntityBrief       `json:"entities"`
}

// EntityBrief is a compact entity representation for context output.
type EntityBrief struct {
	EntityID      string                 `json:"entity_id"`
	Name          string                 `json:"name"`
	State         string                 `json:"state"`
	KeyAttributes map[string]interface{} `json:"key_attributes,omitempty"`
}

// CacheStatsOutput is the output type for cache refresh.
type CacheStatsOutput struct {
	Success     bool   `json:"success"`
	EntityCount int    `json:"entity_count"`
	Timestamp   string `json:"timestamp"`
	CachePath   string `json:"cache_path"`
}

func formatEntityList(v EntityListOutput) string {
	var sb strings.Builder
	for _, e := range v.Entities {
		sb.WriteString(fmt.Sprintf("%-40s %-12s %s\n", e.EntityID, e.State, e.FriendlyName))
	}
	sb.WriteString(fmt.Sprintf("\nTotal: %d entities\n", v.Count))
	return sb.String()
}

func formatEntityState(v EntityStateOutput) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Entity:       %s\n", v.EntityID))
	sb.WriteString(fmt.Sprintf("State:        %s\n", v.State))
	sb.WriteString(fmt.Sprintf("Last Changed: %s\n", v.LastChanged))
	sb.WriteString(fmt.Sprintf("Last Updated: %s\n", v.LastUpdated))
	if len(v.Attributes) > 0 {
		sb.WriteString("Attributes:\n")
		for k, val := range v.Attributes {
			sb.WriteString(fmt.Sprintf("  %s: %v\n", k, val))
		}
	}
	return sb.String()
}

func formatServiceCall(v ServiceCallOutput) string {
	status := "succeeded"
	if !v.Success {
		status = "failed"
	}
	msg := fmt.Sprintf("Service call %s.%s %s", v.Domain, v.Service, status)
	if v.EntityID != "" {
		msg += fmt.Sprintf(" (entity: %s)", v.EntityID)
	}
	return msg + "\n"
}

func formatContext(v ContextOutput) string {
	var sb strings.Builder
	sb.WriteString(v.Summary + "\n\n")
	for _, d := range v.Controllable {
		sb.WriteString(fmt.Sprintf("=== %s (%d) ===\n", d.Domain, d.Count))
		for _, e := range d.Entities {
			sb.WriteString(fmt.Sprintf("  %-35s %-12s %s\n", e.EntityID, e.State, e.Name))
		}
		sb.WriteString("\n")
	}
	if len(v.Other) > 0 {
		sb.WriteString("=== other ===\n")
		for domain, count := range v.Other {
			sb.WriteString(fmt.Sprintf("  %s: %d\n", domain, count))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func formatCacheStats(v CacheStatsOutput) string {
	return fmt.Sprintf("Cache refreshed: %d entities at %s\nPath: %s\n", v.EntityCount, v.Timestamp, v.CachePath)
}
