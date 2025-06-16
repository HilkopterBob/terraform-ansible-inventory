package iohandler

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"reflect"
	"testing"
)

func captureOutput(f func() error) (string, error) {
	r, w, _ := os.Pipe()
	orig := os.Stdout
	os.Stdout = w
	err := f()
	w.Close()
	os.Stdout = orig
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), err
}

func TestOutputJSON(t *testing.T) {
	data := []map[string]interface{}{{"a": "b"}}
	out, err := captureOutput(func() error { return Output(data, "json") })
	if err != nil {
		t.Fatalf("Output returned error: %v", err)
	}
	var got []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(got, data) {
		t.Fatalf("json output mismatch: %v vs %v", got, data)
	}
}

func TestOutputINI(t *testing.T) {
	data := []map[string]interface{}{{"a": "b"}}
	out, err := captureOutput(func() error { return Output(data, "ini") })
	if err != nil {
		t.Fatalf("Output returned error: %v", err)
	}
	expected := "[host0]\na = b\n\n"
	if out != expected {
		t.Fatalf("ini output mismatch:\n%s", out)
	}
}

func TestOutputTXT(t *testing.T) {
	data := []map[string]interface{}{{"a": "b"}}
	out, err := captureOutput(func() error { return Output(data, "txt") })
	if err != nil {
		t.Fatalf("Output returned error: %v", err)
	}
	expected := "Host 0:\n  a: b\n\n"
	if out != expected {
		t.Fatalf("txt output mismatch:\n%s", out)
	}
}

func TestOutputAnsible(t *testing.T) {
	obj := map[string]interface{}{
		"values": map[string]interface{}{
			"name":      "host1",
			"variables": map[string]interface{}{"ip": "1.2.3.4/24"},
		},
	}
	data := []map[string]interface{}{obj}
	out, err := captureOutput(func() error { return Output(data, "ansible") })
	if err != nil {
		t.Fatalf("Output returned error: %v", err)
	}
	expected := "host1 ansible_host=1.2.3.4\n"
	if out != expected {
		t.Fatalf("ansible output mismatch:\n%s", out)
	}
}

func TestOutputUnknown(t *testing.T) {
	err := Output(nil, "bogus")
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestOutputAnsibleInventoryError(t *testing.T) {
	data := []map[string]interface{}{{"a": "b"}}
	_, err := captureOutput(func() error {
		return OutputAnsibleInventory(data, "missing", "ip")
	})
	if err == nil {
		t.Fatal("expected error from OutputAnsibleInventory")
	}
}

func TestLookupDotPath(t *testing.T) {
	obj := map[string]interface{}{
		"a": map[string]interface{}{"b": map[string]interface{}{"c": "val"}},
	}
	v, err := lookupDotPath(obj, "a.b.c")
	if err != nil || v != "val" {
		t.Fatalf("unexpected result %v, %v", v, err)
	}

	if _, err := lookupDotPath(obj, "a.b"); err == nil {
		t.Fatal("expected error for non-string value")
	}
	if _, err := lookupDotPath(obj, "a.x.c"); err == nil {
		t.Fatal("expected error for missing key")
	}
	if _, err := lookupDotPath(obj, "a.b.c.d"); err == nil {
		t.Fatal("expected error for path too deep")
	}
}
