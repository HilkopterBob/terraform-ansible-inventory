package parser

import (
	"github.com/buger/jsonparser"
)

func ExtractAnsibleHosts(data []byte) []map[string]interface{} {
	var results []map[string]interface{}
	stack := [][]byte{data}

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		jsonparser.ObjectEach(current, func(_, value []byte, dataType jsonparser.ValueType, _ int) error {
			switch dataType {
			case jsonparser.Object:
				t, err := jsonparser.GetString(value, "type")
				if err == nil && t == "ansible_host" {
					var obj map[string]interface{}
					if err := jsonparser.Unmarshal(value, &obj); err == nil {
						results = append(results, obj)
					}
				}
				stack = append(stack, value)
			case jsonparser.Array:
				jsonparser.ArrayEach(value, func(elem []byte, dt jsonparser.ValueType, _, _ int, _ error) {
					if dt == jsonparser.Object || dt == jsonparser.Array {
						stack = append(stack, elem)
					}
				})
			}
			return nil
		})
	}
	return results
}
