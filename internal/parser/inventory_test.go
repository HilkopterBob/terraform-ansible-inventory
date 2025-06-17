package parser

import (
	"bytes"
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

func TestParseGroupChildrenParents(t *testing.T) {
	state := map[string]any{
		"values": map[string]any{
			"root_module": map[string]any{
				"resources": []any{
					map[string]any{
						"type": "ansible_group",
						"values": map[string]any{
							"name":      "parent",
							"variables": map[string]string{"role": "p"},
						},
					},
					map[string]any{
						"type": "ansible_group",
						"values": map[string]any{
							"name":     "child",
							"parents":  []any{"parent"},
							"hosts":    []any{"h1"},
							"children": []any{"grand"},
						},
					},
					map[string]any{
						"type": "ansible_group",
						"values": map[string]any{
							"name":  "grand",
							"hosts": []any{"h2"},
						},
					},
					map[string]any{
						"type": "ansible_host",
						"values": map[string]any{
							"name": "h1",
						},
					},
				},
			},
		},
	}
	buf, _ := json.Marshal(state)
	inv := ParseInventory(buf)
	c, ok := inv.Groups["child"]
	if !ok {
		t.Fatalf("child group missing")
	}
	foundParent := false
	for _, p := range c.Parents {
		if p == "parent" {
			foundParent = true
		}
	}
	foundChild := false
	for _, ch := range c.Children {
		if ch == "grand" {
			foundChild = true
		}
	}
	if !foundParent || !foundChild {
		t.Fatalf("child relations missing: %#v", c)
	}
	if _, ok := inv.Hosts["h2"]; !ok {
		t.Fatalf("host h2 from grand child missing")
  }
}
    
func TestParseInventoryReader(t *testing.T) {
	state := map[string]any{
		"type": "ansible_inventory",
		"values": map[string]any{
			"variables": map[string]string{"env": "test"},
		},
	}
	buf, _ := json.Marshal(state)
	inv := ParseInventoryReader(bytes.NewReader(buf))
	if inv.Vars["env"] != "test" {
		t.Fatalf("inventory vars not parsed")
	}
}
