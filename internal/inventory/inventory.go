package inventory

// Inventory holds hosts and groups parsed from Terraform state.
// It loosely mirrors the capabilities of the ansible/ansible provider.
type Inventory struct {
	Hosts  map[string]*Host
	Groups map[string]*Group
	Vars   map[string]string
}

// AddVars merges the provided variables with any existing inventory level
// variables.
func (inv *Inventory) AddVars(v map[string]string) {
	if inv.Vars == nil {
		inv.Vars = make(map[string]string)
	}
	for k, val := range v {
		inv.Vars[k] = val
	}
}

type Host struct {
	Name      string
	Variables map[string]string
	Groups    []string
	Enabled   bool
	Metadata  map[string]string
}

type Group struct {
	Name      string
	Variables map[string]string
	Children  []string
	Hosts     []string
	Parents   []string
}

// New creates an empty Inventory structure.
func New() *Inventory {
	return &Inventory{
		Hosts:  make(map[string]*Host),
		Groups: make(map[string]*Group),
		Vars:   make(map[string]string),
	}
}

// AddHost adds or updates a host.
func (inv *Inventory) AddHost(h *Host) {
	if existing, ok := inv.Hosts[h.Name]; ok {
		// merge variables and groups
		for k, v := range h.Variables {
			existing.Variables[k] = v
		}
		for _, g := range h.Groups {
			if !contains(existing.Groups, g) {
				existing.Groups = append(existing.Groups, g)
			}
		}
		if h.Metadata != nil {
			if existing.Metadata == nil {
				existing.Metadata = make(map[string]string)
			}
			for k, v := range h.Metadata {
				existing.Metadata[k] = v
			}
		}
		if h.Enabled {
			existing.Enabled = h.Enabled
		}
	} else {
		if h.Variables == nil {
			h.Variables = make(map[string]string)
		}
		if h.Metadata == nil {
			h.Metadata = make(map[string]string)
		}
		inv.Hosts[h.Name] = h
	}
	// ensure groups exist
	for _, g := range h.Groups {
		inv.ensureGroup(g)
		grp := inv.Groups[g]
		if !contains(grp.Hosts, h.Name) {
			grp.Hosts = append(grp.Hosts, h.Name)
		}
	}
}

// AddGroup adds or updates a group.
func (inv *Inventory) AddGroup(g *Group) {
	grp := inv.ensureGroup(g.Name)
	for k, v := range g.Variables {
		grp.Variables[k] = v
	}
	for _, child := range g.Children {
		if !contains(grp.Children, child) {
			grp.Children = append(grp.Children, child)
		}
	}
	for _, child := range g.Children {
		inv.ensureGroup(child)
	}
	for _, p := range g.Parents {
		if !contains(grp.Parents, p) {
			grp.Parents = append(grp.Parents, p)
		}
	}
	for _, hname := range g.Hosts {
		if !contains(grp.Hosts, hname) {
			grp.Hosts = append(grp.Hosts, hname)
		}
	}
	for _, h := range g.Hosts {
		if host, ok := inv.Hosts[h]; ok {
			if !contains(host.Groups, g.Name) {
				host.Groups = append(host.Groups, g.Name)
			}
		} else {
			inv.AddHost(&Host{Name: h, Groups: []string{g.Name}})
		}
	}
}

func (inv *Inventory) ensureGroup(name string) *Group {
	if g, ok := inv.Groups[name]; ok {
		return g
	}
	g := &Group{
		Name:      name,
		Variables: make(map[string]string),
	}
	inv.Groups[name] = g
	return g
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func copyMap(src map[string]string) map[string]string {
	if src == nil {
		return nil
	}
	m := make(map[string]string, len(src))
	for k, v := range src {
		m[k] = v
	}
	return m
}

// CopyFiltered creates a new inventory containing only the specified hosts and
// groups. Empty slices mean no filtering on that dimension.
func (inv *Inventory) CopyFiltered(hosts, groups []string) *Inventory {
	hostSet := make(map[string]bool)
	for _, h := range hosts {
		hostSet[h] = true
	}
	groupSet := make(map[string]bool)
	for _, g := range groups {
		groupSet[g] = true
	}

	out := New()
	out.AddVars(inv.Vars)

	for name, h := range inv.Hosts {
		if len(hostSet) > 0 && !hostSet[name] {
			continue
		}
		if len(groupSet) > 0 {
			ok := false
			for _, g := range h.Groups {
				if groupSet[g] {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}
		nh := &Host{
			Name:      h.Name,
			Variables: copyMap(h.Variables),
			Metadata:  copyMap(h.Metadata),
			Groups:    append([]string(nil), h.Groups...),
			Enabled:   h.Enabled,
		}
		out.AddHost(nh)
	}

	for name, g := range inv.Groups {
		if len(groupSet) > 0 && !groupSet[name] {
			continue
		}
		ng := &Group{
			Name:      g.Name,
			Variables: copyMap(g.Variables),
			Children:  append([]string(nil), g.Children...),
			Hosts:     append([]string(nil), g.Hosts...),
			Parents:   append([]string(nil), g.Parents...),
		}
		out.AddGroup(ng)
	}

	return out
}
