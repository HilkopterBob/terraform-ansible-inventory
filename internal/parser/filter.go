package parser

import (
	"encoding/json"

	"github.com/buger/jsonparser"
)

// ExtractAnsibleHosts finds any JSON object at any depth with "type":"ansible_host".
// It traverses both objects and arrays in a stack-based DFS, and always returns a
// non-nil slice (so even if nothing matches, you get "[]" not "null").
func ExtractAnsibleHosts(data []byte) []map[string]interface{} {
	// Pre-allocate empty slice so JSON-encoding yields [] not null
	results := make([]map[string]interface{}, 0, 4)

	// Stack for DFS: start by pushing the root blob
	stack := [][]byte{data}

	for len(stack) > 0 {
		// Pop
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// --- 1) Check this node itself for type=="ansible_host"
		if t, err := jsonparser.GetString(current, "type"); err == nil && t == "ansible_host" {
			var obj map[string]interface{}
			if err := json.Unmarshal(current, &obj); err == nil {
				results = append(results, obj)
			}
		}

		// --- 2) Recurse into any child *objects*
		jsonparser.ObjectEach(current, func(_ []byte, val []byte, dt jsonparser.ValueType, _ int) error {
			// If the field value is an object or an array, push it
			if dt == jsonparser.Object || dt == jsonparser.Array {
				stack = append(stack, val)
			}
			return nil
		})

		// --- 3) Recurse into any child *array* elements
		jsonparser.ArrayEach(current, func(val []byte, dt jsonparser.ValueType, _ int, err error) {
			if err != nil {
				return
			}
			if dt == jsonparser.Object || dt == jsonparser.Array {
				stack = append(stack, val)
			}
		})
	}

	return results
}
