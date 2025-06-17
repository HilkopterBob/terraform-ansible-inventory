package iohandler

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/HilkopterBob/terraform-ansible-inventory/internal/inventory"
	"gopkg.in/yaml.v3"
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
	var data map[string]any
	if err := yaml.Unmarshal([]byte(out), &data); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	all, ok := data["all"].(map[string]any)
	if !ok {
		t.Fatalf("missing all group")
	}
	if hosts, ok := all["hosts"].(map[string]any); ok {
		if _, dup := hosts["test1"]; dup {
			t.Fatalf("host present in all.hosts")
		}
	}
	children, ok := all["children"].(map[string]any)
	if !ok {
		t.Fatalf("missing children")
	}
	web, ok := children["web"].(map[string]any)
	if !ok {
		t.Fatalf("web group missing")
	}
	wh, ok := web["hosts"].(map[string]any)
	if !ok || wh["test1"] == nil {
		t.Fatalf("test1 missing from web group")
	}
	hv, _ := wh["test1"].(map[string]any)
	if hv["ansible_host"] != "192.168.1.10" {
		t.Fatalf("ansible_host missing in yaml")
	}
	vars, _ := all["vars"].(map[string]any)
	if vars["env"] != "test" {
		t.Fatalf("inventory vars missing")
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
	idxHost := strings.Index(out, "test1 ansible_host")
	idxGrp := strings.Index(out, "[web]")
	if idxHost != -1 && idxGrp != -1 && idxHost < idxGrp {
		t.Fatalf("host appears in [all] section:\n%s", out)
	}
}

func TestYAMLvsINIParity(t *testing.T) {
	inv := invFixture()
	yamlOut, err := captureOutput(func() error { return OutputInventory(inv, "yaml") })
	if err != nil {
		t.Fatalf("yaml output error: %v", err)
	}
	iniOut, err := captureOutput(func() error { return OutputInventory(inv, "ini") })
	if err != nil {
		t.Fatalf("ini output error: %v", err)
	}

	var data map[string]any
	if err := yaml.Unmarshal([]byte(yamlOut), &data); err != nil {
		t.Fatalf("unmarshal yaml: %v", err)
	}
	yHost := data["all"].(map[string]any)["children"].(map[string]any)["web"].(map[string]any)["hosts"].(map[string]any)["test1"].(map[string]any)["ansible_host"]
	if yHost != "192.168.1.10" {
		t.Fatalf("unexpected yaml host ip: %v", yHost)
	}
	if !strings.Contains(iniOut, "test1 ansible_host=192.168.1.10") {
		t.Fatalf("ini output missing host ip:\n%s", iniOut)
	}
}
func TestINIYAMLJSONParity(t *testing.T) {
	inv := invFixture()
	yamlOut, err := captureOutput(func() error { return OutputInventory(inv, "yaml") })
	if err != nil {
		t.Fatalf("yaml output error: %v", err)
	}
	iniOut, err := captureOutput(func() error { return OutputInventory(inv, "ini") })
	if err != nil {
		t.Fatalf("ini output error: %v", err)
	}
	jsonOut, err := captureOutput(func() error { return OutputInventory(inv, "json") })
	if err != nil {
		t.Fatalf("json output error: %v", err)
	}

	// parse YAML output
	var ydata map[string]any
	if err := yaml.Unmarshal([]byte(yamlOut), &ydata); err != nil {
		t.Fatalf("unmarshal yaml: %v", err)
	}
	web := ydata["all"].(map[string]any)["children"].(map[string]any)["web"].(map[string]any)
	yhost := web["hosts"].(map[string]any)["test1"].(map[string]any)
	yIP := yhost["ansible_host"].(string)
	yOS := yhost["os"].(string)
	yInvVar := ydata["all"].(map[string]any)["vars"].(map[string]any)["env"].(string)
	yGrpVar := web["vars"].(map[string]any)["tier"].(string)

	// parse JSON output
	var jinv inventory.Inventory
	if err := json.Unmarshal([]byte(jsonOut), &jinv); err != nil {
		t.Fatalf("unmarshal json: %v", err)
	}
	jh := jinv.Hosts["test1"]
	jIP := strings.SplitN(jh.Variables["ip"], "/", 2)[0]
	jOS := jh.Variables["os"]
	jInvVar := jinv.Vars["env"]
	jGrpVar := jinv.Groups["web"].Variables["tier"]

	if yIP != jIP || yOS != jOS {
		t.Fatalf("host variables mismatch yaml vs json")
	}
	if yInvVar != jInvVar {
		t.Fatalf("inventory vars mismatch yaml vs json")
	}
	if yGrpVar != jGrpVar {
		t.Fatalf("group vars mismatch yaml vs json")
	}

	expectedHost := "test1 ansible_host=" + jIP + " os=" + jOS
	if !strings.Contains(iniOut, expectedHost) {
		t.Fatalf("ini output missing host line: %s", expectedHost)
	}
	if !strings.Contains(iniOut, "[all:vars]") || !strings.Contains(iniOut, "env="+jInvVar) {
		t.Fatalf("ini output missing inventory vars")
	}
	if !strings.Contains(iniOut, "[web:vars]") || !strings.Contains(iniOut, "tier="+jGrpVar) {
		t.Fatalf("ini output missing group vars")
	}
}

func TestOutputInventoryUnknownFormat(t *testing.T) {
	inv := inventory.New()
	err := OutputInventory(inv, "bogus")
	if err == nil {
		t.Fatal("expected error for unknown inventory format")
	}
}

func TestStripCIDRHelper(t *testing.T) {
	if stripCIDR("1.2.3.4/32") != "1.2.3.4" {
		t.Fatalf("stripCIDR failed")
	}
	if stripCIDR("1.2.3.4") != "1.2.3.4" {
		t.Fatalf("stripCIDR modified plain ip")
	}
}
