package inventory

// Inventory holds hosts and groups parsed from Terraform state.
// It loosely mirrors the capabilities of the ansible/ansible provider.
type Inventory struct {
	Hosts  map[string]*Host
	Groups map[string]*Group
}

type Host struct {
	Name      string
	Variables map[string]string
	Groups    []string
}

type Group struct {
	Name      string
	Variables map[string]string
	Children  []string
	Hosts     []string
}

// New creates an empty Inventory structure.
func New() *Inventory {
	return &Inventory{
		Hosts:  make(map[string]*Host),
		Groups: make(map[string]*Group),
	}
}

// AddHost adds or updates a host.
func (inv *Inventory) AddHost(h *Host) {
	if existing, ok := inv.Hosts[h.Name]; ok {
		// merge variables and groups
		for k, v := range h.Variables {
			existing.Variables[k] = v
		}
		existing.Groups = append(existing.Groups, h.Groups...)
	} else {
		if h.Variables == nil {
			h.Variables = make(map[string]string)
		}
		inv.Hosts[h.Name] = h
	}
	// ensure groups exist
	for _, g := range h.Groups {
		inv.ensureGroup(g)
		grp := inv.Groups[g]
		grp.Hosts = append(grp.Hosts, h.Name)
	}
}

// AddGroup adds or updates a group.
func (inv *Inventory) AddGroup(g *Group) {
	grp := inv.ensureGroup(g.Name)
	for k, v := range g.Variables {
		grp.Variables[k] = v
	}
	grp.Children = append(grp.Children, g.Children...)
	for _, child := range g.Children {
		inv.ensureGroup(child)
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
