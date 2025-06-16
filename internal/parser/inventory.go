package parser

import (
	"encoding/json"

	"github.com/HilkopterBob/terraform-ansible-inventory/internal/inventory"
	"github.com/buger/jsonparser"
)

// TODO: Expand parsing logic to cover all resources exposed by the
// ansible/ansible Terraform provider. This currently only understands
// ansible_host and ansible_group, but the provider includes additional
// resources such as group membership and inventory level variables that
// should be reflected in the Inventory structure.

// ParseInventory walks the Terraform state JSON and extracts all ansible_host
// and ansible_group resources, returning a structured Inventory.
func ParseInventory(data []byte) *inventory.Inventory {
	inv := inventory.New()
	stack := [][]byte{data}

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Determine if this object is a resource we care about
		if t, err := jsonparser.GetString(current, "type"); err == nil {
			// TODO: Handle additional resource types emitted by the
			// Terraform provider once they are added.
			switch t {
			case "ansible_host":
				var tmp struct {
					Values struct {
						Name      string            `json:"name"`
						Groups    []string          `json:"groups"`
						Variables map[string]string `json:"variables"`
					} `json:"values"`
				}
				if err := json.Unmarshal(current, &tmp); err == nil {
					inv.AddHost(&inventory.Host{
						Name:      tmp.Values.Name,
						Groups:    tmp.Values.Groups,
						Variables: tmp.Values.Variables,
					})
				}
			case "ansible_group":
				var tmp struct {
					Values struct {
						Name      string            `json:"name"`
						Children  []string          `json:"children"`
						Variables map[string]string `json:"variables"`
					} `json:"values"`
				}
				if err := json.Unmarshal(current, &tmp); err == nil {
					inv.AddGroup(&inventory.Group{
						Name:      tmp.Values.Name,
						Children:  tmp.Values.Children,
						Variables: tmp.Values.Variables,
					})
				}
			}
		}

		// Recurse into child objects and arrays
		jsonparser.ObjectEach(current, func(_ []byte, val []byte, dt jsonparser.ValueType, _ int) error {
			if dt == jsonparser.Object || dt == jsonparser.Array {
				stack = append(stack, val)
			}
			return nil
		})
		jsonparser.ArrayEach(current, func(val []byte, dt jsonparser.ValueType, _ int, err error) {
			if err != nil {
				return
			}
			if dt == jsonparser.Object || dt == jsonparser.Array {
				stack = append(stack, val)
			}
		})
	}

	return inv
}
