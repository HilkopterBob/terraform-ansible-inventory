package iohandler

import (
	"strings"
	"testing"

	"github.com/HilkopterBob/terraform-ansible-inventory/internal/inventory"
)

func TestINIFilterNoDuplicates(t *testing.T) {
	inv := inventory.New()
	inv.AddHost(&inventory.Host{Name: "h1", Groups: []string{"web"}, Variables: map[string]string{"ip": "1.2.3.4"}})
	inv.AddGroup(&inventory.Group{Name: "web", Hosts: []string{"h1"}})

	inv = inv.CopyFiltered([]string{"h1"}, nil)
	out, err := captureOutput(func() error { return OutputInventory(inv, "ini") })
	if err != nil {
		t.Fatalf("ini output error: %v", err)
	}
	count := strings.Count(out, "h1 ansible_host=1.2.3.4")
	if count != 1 {
		t.Fatalf("expected 1 host entry, got %d: \n%s", count, out)
	}
}
