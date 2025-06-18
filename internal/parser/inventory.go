package parser

import (
	"bytes"
	"io"

	"github.com/bcicen/jstream"

	"github.com/HilkopterBob/terraform-ansible-inventory/internal/inventory"
)

// ParseInventory walks the Terraform state JSON and extracts ansible_* resources
// to build an inventory compatible with the ansible/ansible provider.

// ParseInventory walks the Terraform state JSON and extracts all ansible_host
// and ansible_group resources, returning a structured Inventory.
func ParseInventory(data []byte) *inventory.Inventory {
	return ParseInventoryReader(bytes.NewReader(data))
}

// ParseInventoryReader streams Terraform state JSON from r and extracts
// ansible_* resources to build an inventory compatible with the
// ansible/ansible provider.
func ParseInventoryReader(r io.Reader) *inventory.Inventory {
	inv := inventory.New()
	dec := jstream.NewDecoder(r, -1)

	for mv := range dec.Stream() {
		if mv.ValueType != jstream.Object {
			continue
		}
		obj, ok := mv.Value.(map[string]interface{})
		if !ok {
			continue
		}
		t, _ := obj["type"].(string)
		values, _ := obj["values"].(map[string]interface{})

		switch t {
		case "ansible_host":
			h := &inventory.Host{
				Name:      getString(values["name"]),
				Groups:    toStringSlice(values["groups"]),
				Variables: toStringMap(values["variables"]),
				Metadata:  toStringMap(values["metadata"]),
				Enabled:   getBool(values["enabled"]),
			}
			inv.AddHost(h)
		case "ansible_group":
			g := &inventory.Group{
				Name:      getString(values["name"]),
				Children:  toStringSlice(values["children"]),
				Variables: toStringMap(values["variables"]),
				Hosts:     toStringSlice(values["hosts"]),
				Parents:   toStringSlice(values["parents"]),
			}
			inv.AddGroup(g)
		case "ansible_inventory":
			inv.AddVars(toStringMap(values["variables"]))
		}
	}

	return inv
}

func getString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func toStringSlice(v interface{}) []string {
	arr, ok := v.([]interface{})
	if !ok {
		return nil
	}
	out := make([]string, 0, len(arr))
	for _, x := range arr {
		if s, ok := x.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func toStringMap(v interface{}) map[string]string {
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, val := range m {
		if s, ok := val.(string); ok {
			out[k] = s
		}
	}
	return out
}

func getBool(v interface{}) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}
