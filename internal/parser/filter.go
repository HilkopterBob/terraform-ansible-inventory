package parser

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/bcicen/jstream"
)

// ExtractAnsibleHosts finds any JSON object at any depth with
// "type":"ansible_host". The returned slice is never nil even when no
// matches are found.
func ExtractAnsibleHosts(data []byte) []map[string]interface{} {
	return ExtractAnsibleHostsReader(bytes.NewReader(data))
}

// ExtractAnsibleHostsReader streams JSON from r and returns all objects where
// "type" == "ansible_host". The returned slice is never nil.
func ExtractAnsibleHostsReader(r io.Reader) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, 4)
	dec := jstream.NewDecoder(r, -1)

	for mv := range dec.Stream() {
		if mv.ValueType != jstream.Object {
			continue
		}
		obj, ok := mv.Value.(map[string]interface{})
		if !ok {
			continue
		}
		if t, _ := obj["type"].(string); t == "ansible_host" {
			// copy so callers can modify without affecting parser state
			buf, _ := json.Marshal(obj)
			var out map[string]interface{}
			if err := json.Unmarshal(buf, &out); err == nil {
				results = append(results, out)
			}
		}
	}

	return results
}
