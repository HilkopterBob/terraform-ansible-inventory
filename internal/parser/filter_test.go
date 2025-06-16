package parser

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestExtractAnsibleHosts(t *testing.T) {
	sample := map[string]interface{}{
		"values": map[string]interface{}{
			"root_module": map[string]interface{}{
				"child_modules": []interface{}{
					map[string]interface{}{
						"resources": []interface{}{
							map[string]interface{}{
								"type": "ansible_host",
								"values": map[string]interface{}{
									"name":      "host1",
									"variables": map[string]interface{}{"ip": "10.0.0.1/24"},
								},
							},
							map[string]interface{}{"type": "other"},
							map[string]interface{}{
								"type": "ansible_host",
								"values": map[string]interface{}{
									"name":      "host2",
									"variables": map[string]interface{}{"ip": "10.0.0.2"},
								},
							},
						},
						"nested": map[string]interface{}{
							"type": "ansible_host",
							"values": map[string]interface{}{
								"name":      "host3",
								"variables": map[string]interface{}{"ip": "10.0.0.3"},
							},
						},
					},
				},
			},
		},
	}
	buf, err := json.Marshal(sample)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	got := ExtractAnsibleHosts(buf)
	if len(got) != 3 {
		t.Fatalf("expected 3 hosts, got %d", len(got))
	}

	wantNames := map[string]bool{"host1": true, "host2": true, "host3": true}
	for _, obj := range got {
		if obj["type"] != "ansible_host" {
			t.Errorf("unexpected type %v", obj["type"])
		}
		values, ok := obj["values"].(map[string]interface{})
		if !ok {
			t.Fatalf("values not a map: %T", obj["values"])
		}
		name, _ := values["name"].(string)
		if !wantNames[name] {
			t.Errorf("unexpected host name %s", name)
		}
		delete(wantNames, name)
	}
	if len(wantNames) != 0 {
		t.Errorf("missing hosts: %v", reflect.ValueOf(wantNames).MapKeys())
	}
}

func TestExtractAnsibleHostsNone(t *testing.T) {
	data := []byte(`{"foo":1}`)
	got := ExtractAnsibleHosts(data)
	if got == nil {
		t.Fatal("expected empty slice, got nil")
	}
	if len(got) != 0 {
		t.Fatalf("expected 0 hosts, got %d", len(got))
	}
}
