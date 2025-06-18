package inventory

import "testing"

func TestAddVars(t *testing.T) {
	inv := New()
	inv.AddVars(map[string]string{"a": "1"})
	inv.AddVars(map[string]string{"b": "2", "a": "x"})
	if inv.Vars["a"] != "x" || inv.Vars["b"] != "2" {
		t.Fatalf("vars not merged correctly: %#v", inv.Vars)
	}
}

func TestAddHostNewAndMerge(t *testing.T) {
	inv := New()
	inv.AddHost(&Host{Name: "h1", Groups: []string{"g1"}, Variables: map[string]string{"ip": "1"}})
	inv.AddHost(&Host{Name: "h1", Groups: []string{"g2"}, Variables: map[string]string{"os": "linux"}, Metadata: map[string]string{"role": "db"}, Enabled: true})

	h, ok := inv.Hosts["h1"]
	if !ok {
		t.Fatalf("host missing")
	}
	if h.Variables["ip"] != "1" || h.Variables["os"] != "linux" {
		t.Fatalf("variables not merged: %#v", h.Variables)
	}
	if len(h.Groups) != 2 || !(h.Groups[0] == "g1" && h.Groups[1] == "g2" || h.Groups[0] == "g2" && h.Groups[1] == "g1") {
		t.Fatalf("groups not merged: %#v", h.Groups)
	}
	if h.Metadata["role"] != "db" {
		t.Fatalf("metadata not merged: %#v", h.Metadata)
	}
	if !h.Enabled {
		t.Fatalf("enabled flag not set")
	}
	if g := inv.Groups["g2"]; g == nil || len(g.Hosts) == 0 || g.Hosts[0] != "h1" {
		t.Fatalf("group g2 not updated with host")
	}
}

func TestAddGroupNewAndMerge(t *testing.T) {
	inv := New()
	inv.AddGroup(&Group{Name: "web", Variables: map[string]string{"tier": "fe"}, Hosts: []string{"h1"}, Children: []string{"child"}, Parents: []string{"parent"}})
	inv.AddGroup(&Group{Name: "web", Variables: map[string]string{"env": "prod"}, Hosts: []string{"h2"}, Children: []string{"child2"}})

	g, ok := inv.Groups["web"]
	if !ok {
		t.Fatalf("group missing")
	}
	if g.Variables["tier"] != "fe" || g.Variables["env"] != "prod" {
		t.Fatalf("variables not merged: %#v", g.Variables)
	}
	if !(contains(g.Hosts, "h1") && contains(g.Hosts, "h2")) {
		t.Fatalf("hosts not merged: %#v", g.Hosts)
	}
	if len(g.Children) != 2 {
		t.Fatalf("children not merged: %#v", g.Children)
	}
	if len(g.Parents) != 1 || g.Parents[0] != "parent" {
		t.Fatalf("parents not set: %#v", g.Parents)
	}
	if h := inv.Hosts["h2"]; h == nil || (len(h.Groups) == 0 || h.Groups[0] != "web" && h.Groups[1] != "web") {
		t.Fatalf("host h2 not added with group")
	}
	if _, ok := inv.Groups["child"]; !ok {
		t.Fatalf("child group missing")
	}
	if _, ok := inv.Groups["child2"]; !ok {
		t.Fatalf("child2 group missing")
	}
}

func TestCopyFiltered(t *testing.T) {
	inv := New()
	inv.AddVars(map[string]string{"env": "test"})
	inv.AddHost(&Host{Name: "h1", Groups: []string{"web"}})
	inv.AddHost(&Host{Name: "h2", Groups: []string{"db"}})
	inv.AddGroup(&Group{Name: "web", Hosts: []string{"h1"}})
	inv.AddGroup(&Group{Name: "db", Hosts: []string{"h2"}})

	c1 := inv.CopyFiltered([]string{"h1"}, []string{"web"})
	if len(c1.Hosts) != 1 || c1.Hosts["h1"] == nil {
		t.Fatalf("host filter failed")
	}
	if len(c1.Groups) != 1 || c1.Groups["web"] == nil {
		t.Fatalf("group filter mismatch")
	}

	c2 := inv.CopyFiltered(nil, []string{"db"})
	if len(c2.Hosts) != 1 || c2.Hosts["h2"] == nil {
		t.Fatalf("group filter host selection failed")
	}
	if len(c2.Groups) != 1 || c2.Groups["db"] == nil {
		t.Fatalf("group filter failed")
	}
}
func TestAddHostNilMaps(t *testing.T) {
	inv := New()
	inv.AddHost(&Host{Name: "x1"})
	h, ok := inv.Hosts["x1"]
	if !ok {
		t.Fatalf("host not created")
	}
	if h.Variables == nil || h.Metadata == nil {
		t.Fatalf("nil maps after AddHost: vars=%v meta=%v", h.Variables, h.Metadata)
	}
	inv.AddHost(&Host{Name: "x1", Variables: map[string]string{"os": "linux"}})
	if h.Variables["os"] != "linux" {
		t.Fatalf("variable merge failed: %#v", h.Variables)
	}
}

func TestAddGroupCreatesHosts(t *testing.T) {
	inv := New()
	inv.AddGroup(&Group{Name: "db", Hosts: []string{"db1"}, Children: []string{"child"}})
	if h := inv.Hosts["db1"]; h == nil || !contains(h.Groups, "db") {
		t.Fatalf("host db1 not created with group")
	}
	if _, ok := inv.Groups["child"]; !ok {
		t.Fatalf("child group not created")
	}
}
func TestMergeGroupHostDuplicatesEdge(t *testing.T) {
	inv := New()
	inv.AddGroup(&Group{Name: "dup", Hosts: []string{"h1", "h1"}})
	if len(inv.Groups["dup"].Hosts) != 1 {
		t.Fatalf("expected deduplicated host list, got %v", inv.Groups["dup"].Hosts)
	}
	inv.AddHost(&Host{Name: "h1", Groups: []string{"dup", "dup"}})
	if len(inv.Hosts["h1"].Groups) != 1 {
		t.Fatalf("expected deduplicated group list on host, got %v", inv.Hosts["h1"].Groups)
	}
}

func TestAddHostInvalidName(t *testing.T) {
	inv := New()
	inv.AddHost(&Host{Name: "", Groups: []string{"web"}})
	if _, ok := inv.Hosts[""]; !ok {
		t.Fatalf("host with empty name not present")
	}
	inv.AddHost(&Host{Name: "", Groups: []string{"web"}})
	if len(inv.Groups["web"].Hosts) != 1 {
		t.Fatalf("expected deduped empty host in group, got %v", inv.Groups["web"].Hosts)
	}
}
