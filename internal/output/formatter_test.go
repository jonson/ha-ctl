package output

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestJSONFormatterEntityList(t *testing.T) {
	f := &JSONFormatter{}
	data := EntityListOutput{
		Entities: []EntityItem{
			{EntityID: "light.kitchen", State: "on", FriendlyName: "Kitchen", Domain: "light"},
		},
		Count: 1,
	}
	result, err := f.Format(data)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Should be valid JSON
	var parsed map[string]any
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if !strings.Contains(result, "light.kitchen") {
		t.Error("expected output to contain light.kitchen")
	}
}

func TestTextFormatterEntityList(t *testing.T) {
	f := &TextFormatter{}
	data := EntityListOutput{
		Entities: []EntityItem{
			{EntityID: "light.kitchen", State: "on", FriendlyName: "Kitchen", Domain: "light"},
			{EntityID: "switch.garage", State: "off", FriendlyName: "Garage", Domain: "switch"},
		},
		Count: 2,
	}
	result, err := f.Format(data)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	if !strings.Contains(result, "light.kitchen") {
		t.Error("expected output to contain light.kitchen")
	}
	if !strings.Contains(result, "Total: 2") {
		t.Error("expected output to contain total count")
	}
}

func TestJSONFormatterServiceCall(t *testing.T) {
	f := &JSONFormatter{}
	data := ServiceCallOutput{
		Success:  true,
		Domain:   "light",
		Service:  "turn_on",
		EntityID: "light.kitchen",
	}
	result, err := f.Format(data)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if parsed["success"] != true {
		t.Error("expected success=true")
	}
}

func TestJSONFormatterContext(t *testing.T) {
	f := &JSONFormatter{}
	data := ContextOutput{
		Summary: "Home has 10 entities across 3 domains",
		Controllable: []DomainSummary{
			{
				Domain: "light",
				Count:  2,
				Entities: []EntityBrief{
					{EntityID: "light.kitchen", Name: "Kitchen", State: "on"},
					{EntityID: "light.bedroom", Name: "Bedroom", State: "off"},
				},
			},
		},
		Other: map[string]int{
			"sensor": 5,
			"update": 3,
		},
	}
	result, err := f.Format(data)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if !strings.Contains(result, "controllable") {
		t.Error("expected output to contain 'controllable'")
	}
	if !strings.Contains(result, "other") {
		t.Error("expected output to contain 'other'")
	}
	if !strings.Contains(result, "light.kitchen") {
		t.Error("expected output to contain light.kitchen")
	}
}

func TestTextFormatterContext(t *testing.T) {
	f := &TextFormatter{}
	data := ContextOutput{
		Summary: "Home has 10 entities across 3 domains",
		Controllable: []DomainSummary{
			{
				Domain: "light",
				Count:  1,
				Entities: []EntityBrief{
					{EntityID: "light.kitchen", Name: "Kitchen", State: "on"},
				},
			},
		},
		Other: map[string]int{
			"sensor": 5,
		},
	}
	result, err := f.Format(data)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	if !strings.Contains(result, "light.kitchen") {
		t.Error("expected output to contain light.kitchen")
	}
	if !strings.Contains(result, "=== other ===") {
		t.Error("expected output to contain '=== other ===' section")
	}
	if !strings.Contains(result, "sensor: 5") {
		t.Error("expected output to contain 'sensor: 5'")
	}
}

func TestNewFormatterJSON(t *testing.T) {
	f := New("json")
	if _, ok := f.(*JSONFormatter); !ok {
		t.Error("expected JSONFormatter for 'json'")
	}
}

func TestNewFormatterText(t *testing.T) {
	f := New("text")
	if _, ok := f.(*TextFormatter); !ok {
		t.Error("expected TextFormatter for 'text'")
	}
}

func TestNewFormatterDefault(t *testing.T) {
	f := New("anything")
	if _, ok := f.(*JSONFormatter); !ok {
		t.Error("expected JSONFormatter as default")
	}
}
