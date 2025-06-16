package iohandler

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/HilkopterBob/terraform-ansible-inventory/internal/inventory"
)

func newTestInventory() *inventory.Inventory {
	inv := inventory.New()
	inv.AddHost(&inventory.Host{Name: "ns0.lxc.tftest.lab.dtbnet.de", Variables: map[string]string{"ansible_host": "10.93.80.200"}, Groups: []string{"bind"}})
	inv.AddHost(&inventory.Host{Name: "ns1.lxc.tftest.lab.dtbnet.de", Variables: map[string]string{"ansible_host": "10.93.80.201"}, Groups: []string{"bind"}})
	inv.AddHost(&inventory.Host{Name: "ns2.lxc.tftest.lab.dtbnet.de", Variables: map[string]string{"ansible_host": "10.93.80.202"}, Groups: []string{"bind"}})
	inv.AddGroup(&inventory.Group{Name: "bind"})
	return inv
}

func TestOutputInventoryYAMLGroups(t *testing.T) {
	inv := newTestInventory()
	out, err := captureOutput(func() error { return OutputInventory(inv, "yaml") })
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	var data map[string]any
	if err := yaml.Unmarshal([]byte(out), &data); err != nil {
		t.Fatalf("unmarshal yaml: %v", err)
	}
	all := data["all"].(map[string]any)
	if hosts, ok := all["hosts"]; ok && len(hosts.(map[string]any)) != 0 {
		t.Fatalf("expected no hosts under all, got %v", hosts)
	}
	bind := all["children"].(map[string]any)["bind"].(map[string]any)
	bh := bind["hosts"].(map[string]any)
	if len(bh) != 3 {
		t.Fatalf("expected 3 bind hosts, got %d", len(bh))
	}
}

func TestOutputInventoryINIGroup(t *testing.T) {
	inv := newTestInventory()
	out, err := captureOutput(func() error { return OutputInventory(inv, "ini") })
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	lines := []string{}
	for _, line := range splitLines(out) {
		if line != "" {
			lines = append(lines, line)
		}
	}
	if len(lines) != 4 || lines[0] != "[bind]" {
		t.Fatalf("unexpected ini output:\n%s", out)
	}
}

func TestOutputInventoryJSONGroup(t *testing.T) {
	inv := newTestInventory()
	out, err := captureOutput(func() error { return OutputInventory(inv, "json") })
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	var data map[string]any
	if err := json.Unmarshal([]byte(out), &data); err != nil {
		t.Fatalf("unmarshal json: %v", err)
	}
	if len(data["Hosts"].(map[string]any)) != 0 {
		t.Fatalf("expected no hosts in Hosts field")
	}
	bind := data["Groups"].(map[string]any)["bind"].(map[string]any)
	bh := bind["Hosts"].([]interface{})
	if len(bh) != 3 {
		t.Fatalf("expected 3 bind hosts, got %d", len(bh))
	}
}

// helper
func splitLines(s string) []string {
	lines := []string{}
	cur := ""
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, cur)
			cur = ""
		} else {
			cur += string(s[i])
		}
	}
	if cur != "" {
		lines = append(lines, cur)
	}
	return lines
}
