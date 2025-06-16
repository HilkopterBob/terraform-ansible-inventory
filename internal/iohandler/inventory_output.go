package iohandler

import (
	"encoding/json"
	"fmt"
	"sort"

	"gopkg.in/yaml.v3"

	"github.com/HilkopterBob/terraform-ansible-inventory/internal/inventory"
)

type groupYAML struct {
	Hosts    map[string]any        `yaml:"hosts,omitempty"`
	Vars     map[string]string     `yaml:"vars,omitempty"`
	Children map[string]*groupYAML `yaml:"children,omitempty"`
}

// OutputInventory dispatches YAML or INI inventory output.
func OutputInventory(inv *inventory.Inventory, format string) error {
	// TODO: Add export formats that mirror Ansible's inventory plugins and
	// ensure parity with the provider's data model.
	switch format {
	case "json":
		return outputJSONInventory(inv)
	case "yaml":
		return outputYAML(inv)
	case "ini":
		return outputINIInventory(inv)
	default:
		return fmt.Errorf("unknown inventory format: %s", format)
	}
}

func outputYAML(inv *inventory.Inventory) error {
	root := &groupYAML{
		Hosts:    make(map[string]any),
		Children: make(map[string]*groupYAML),
	}

	// TODO: Review this structure against Ansible's YAML inventory
	// specification. The current output may not cover all supported
	// attributes like group-level vars or nested hostvars.

	// add hosts
	for _, h := range inv.Hosts {
		hostVars := make(map[string]string)
		for k, v := range h.Variables {
			hostVars[k] = v
		}
		if ip, ok := hostVars["ip"]; ok {
			hostVars["ansible_host"] = stripCIDR(ip)
			delete(hostVars, "ip")
		}
		if len(hostVars) == 0 {
			root.Hosts[h.Name] = struct{}{}
		} else {
			root.Hosts[h.Name] = hostVars
		}
	}

	// prepare groups recursively
	for _, g := range inv.Groups {
		gy := ensureGroupYAML(root, g.Name)
		if len(g.Variables) > 0 {
			if gy.Vars == nil {
				gy.Vars = make(map[string]string)
			}
			for k, v := range g.Variables {
				gy.Vars[k] = v
			}
		}
		for _, child := range g.Children {
			ensureGroupYAML(root, g.Name).Children[child] = ensureGroupYAML(root, child)
		}
		for _, host := range g.Hosts {
			if gy.Hosts == nil {
				gy.Hosts = make(map[string]any)
			}
			gy.Hosts[host] = struct{}{}
		}
	}

	out := map[string]*groupYAML{"all": root}
	enc := yaml.NewEncoder(stdoutWrapper{})
	enc.SetIndent(2)
	if err := enc.Encode(out); err != nil {
		return err
	}
	return enc.Close()
}

type stdoutWrapper struct{}

func (stdoutWrapper) Write(p []byte) (int, error) { return fmt.Print(string(p)) }

func ensureGroupYAML(root *groupYAML, name string) *groupYAML {
	parts := []string{name}
	gy := root
	for _, p := range parts {
		if gy.Children == nil {
			gy.Children = make(map[string]*groupYAML)
		}
		if _, ok := gy.Children[p]; !ok {
			gy.Children[p] = &groupYAML{Children: make(map[string]*groupYAML)}
		}
		gy = gy.Children[p]
	}
	return gy
}

func stripCIDR(ip string) string {
	if idx := len(ip); idx > 0 {
		if pos := index(ip, '/'); pos >= 0 {
			return ip[:pos]
		}
	}
	return ip
}

func index(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}

func outputINIInventory(inv *inventory.Inventory) error {
	var out string

	// TODO: Rework INI generation to produce valid Ansible inventory
	// syntax including group nesting and host variable sections.

	groups := make([]string, 0, len(inv.Groups))
	for name := range inv.Groups {
		groups = append(groups, name)
	}
	sort.Strings(groups)

	// hosts not in any group -> all
	out += "[all]\n"
	for _, h := range inv.Hosts {
		if len(h.Groups) == 0 {
			out += formatHostINI(h) + "\n"
		}
	}
	out += "\n"

	for _, gname := range groups {
		g := inv.Groups[gname]
		if len(g.Hosts) > 0 {
			out += fmt.Sprintf("[%s]\n", gname)
			for _, h := range g.Hosts {
				out += h + "\n"
			}
			out += "\n"
		}
		if len(g.Variables) > 0 {
			out += fmt.Sprintf("[%s:vars]\n", gname)
			for k, v := range g.Variables {
				out += fmt.Sprintf("%s=%s\n", k, v)
			}
			out += "\n"
		}
		if len(g.Children) > 0 {
			out += fmt.Sprintf("[%s:children]\n", gname)
			for _, c := range g.Children {
				out += c + "\n"
			}
			out += "\n"
		}
	}

	_, err := fmt.Print(out)
	return err
}

func formatHostINI(h *inventory.Host) string {
	line := h.Name
	if ip, ok := h.Variables["ip"]; ok {
		line += fmt.Sprintf(" ansible_host=%s", stripCIDR(ip))
	}
	for k, v := range h.Variables {
		if k == "ip" {
			continue
		}
		line += fmt.Sprintf(" %s=%s", k, v)
	}
	// TODO: include all supported host parameters from the provider when
	// formatting the INI entry.
	return line
}

func outputJSONInventory(inv *inventory.Inventory) error {
	enc := json.NewEncoder(stdoutWrapper{})
	enc.SetIndent("", "  ")
	return enc.Encode(inv)
}
