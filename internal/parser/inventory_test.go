package parser

import (
	"encoding/json"
	"testing"
)

func TestParseInventory(t *testing.T) {
	state := map[string]any{
		"values": map[string]any{
			"root_module": map[string]any{
				"resources": []any{
					map[string]any{
						"type": "ansible_inventory",
						"values": map[string]any{
							"variables": map[string]string{"env": "test"},
						},
					},
					map[string]any{
						"type": "ansible_group",
						"values": map[string]any{
							"name":      "web",
							"variables": map[string]string{"tier": "frontend"},
						},
					},
					map[string]any{
						"type": "ansible_host",
						"values": map[string]any{
							"name":      "host1",
							"groups":    []any{"web"},
							"variables": map[string]string{"ip": "10.0.0.1"},
						},
					},
				},
			},
		},
	}
	buf, _ := json.Marshal(state)
	inv := ParseInventory(buf)
	if inv.Vars["env"] != "test" {
		t.Fatalf("inventory vars not parsed")
	}
	if len(inv.Groups) != 1 || inv.Groups["web"].Variables["tier"] != "frontend" {
		t.Fatalf("group parse failed")
	}
	if h, ok := inv.Hosts["host1"]; !ok || h.Variables["ip"] != "10.0.0.1" {
		t.Fatalf("host parse failed")
	}
}
