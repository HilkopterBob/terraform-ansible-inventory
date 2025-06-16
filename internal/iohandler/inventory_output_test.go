package iohandler

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/HilkopterBob/terraform-ansible-inventory/internal/inventory"
)

func invFixture() *inventory.Inventory {
	inv := inventory.New()
	inv.AddVars(map[string]string{"env": "test"})
	inv.AddGroup(&inventory.Group{Name: "web", Variables: map[string]string{"tier": "frontend"}})
	inv.AddHost(&inventory.Host{
		Name:      "test1",
		Groups:    []string{"web"},
		Variables: map[string]string{"ip": "192.168.1.10/24", "os": "linux"},
	})
	return inv
}

func TestOutputInventoryJSON(t *testing.T) {
	inv := invFixture()
	out, err := captureOutput(func() error { return OutputInventory(inv, "json") })
	if err != nil {
		t.Fatalf("json output error: %v", err)
	}
	var chk inventory.Inventory
	if err := json.Unmarshal([]byte(out), &chk); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if chk.Hosts["test1"].Variables["os"] != "linux" {
		t.Fatalf("unexpected json output")
	}
}

func TestOutputInventoryYAML(t *testing.T) {
	inv := invFixture()
	out, err := captureOutput(func() error { return OutputInventory(inv, "yaml") })
	if err != nil {
		t.Fatalf("yaml output error: %v", err)
	}
	if !strings.Contains(out, "test1:") || !strings.Contains(out, "env:") {
		t.Fatalf("unexpected yaml output: %s", out)
	}
}

func TestOutputInventoryINI(t *testing.T) {
	inv := invFixture()
	out, err := captureOutput(func() error { return OutputInventory(inv, "ini") })
	if err != nil {
		t.Fatalf("ini output error: %v", err)
	}
	if !strings.Contains(out, "[web]") || !strings.Contains(out, "ansible_host=192.168.1.10") {
		t.Fatalf("unexpected ini output:\n%s", out)
	}
}
