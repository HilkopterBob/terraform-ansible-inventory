# Terraform Ansible Inventory

## About

`terraform-ansible-inventory` is a small Go CLI for converting Terraform state
into complete Ansible inventories. The tool was created to work with the
`ansible/ansible` Terraform provider which stores hosts, groups and variables in
the state file. RedHat also provides a dynamic inventory plugin that reads
Terraform state, but it cannot handle child modules and therefore misses large
parts of real-world configurations. This CLI parses the state directly and
fully supports nested modules so you can build inventories reliably.

To generate usable state data you **must** manage your inventory with the
`ansible/ansible` provider. Once the provider has written hosts and groups into
your state file, this tool can consume it and output Ansible formatted data.


## Features

- **Streaming JSON parsing** for huge state files without high memory use.
- **Multiple output formats**: `yaml`, `ini`, `json` and native Ansible
  inventory.
- **Understands provider resources**: host variables, group hierarchy and
  inventory level variables from the `ansible/ansible` provider.
- **Child module aware**: traverses nested modules to pick up all resources.
- **Built-in IP/CIDR handling** so exported addresses work directly in Ansible.
- **Clean CLI interface** with automatic `--help` and sensible defaults.
- **No external runtime dependencies** other than the Go binary itself.
- **CI ready**: sample state file and GitHub Actions workflow included.
- Filter output by host or group using `--host` and `--group` flags.

## Installation

Clone and build from source or grab a prebuilt binary from the
[releases page](https://github.com/HilkopterBob/terraform-ansible-inventory/releases):

```bash
# Build from source
git clone https://github.com/HilkopterBob/terraform-ansible-inventory.git
cd terraform-ansible-inventory

go build -o terraform-ansible-inventory ./main.go

# (Optional) install globally
go install github.com/HilkopterBob/terraform-ansible-inventory@latest
```

Download the latest release if you prefer not to build the binary yourself.

## Usage

Run the executable with `--help` to see all options:

```bash
terraform-ansible-inventory --help
```

Common examples assuming your state file is `state.json`:

```bash
# YAML inventory
terraform-ansible-inventory -i state.json -f yaml > inventory.yml

# INI inventory
terraform-ansible-inventory -i state.json -f ini > inventory.ini

# JSON machine readable form
terraform-ansible-inventory -i state.json -f json > inventory.json

# Native Ansible inventory
terraform-ansible-inventory -i state.json -f ansible
```

You can restrict output to specific hosts or groups using the `--host` and
`--group` flags. Multiple values are allowed.

## Advanced configuration

You can override the JSON paths for hostnames and IP addresses if your state
uses custom keys:

```bash
terraform-ansible-inventory \
  --input state.json \
  --format ansible \
  --host-field "values.custom_hostname" \
  --ip-field "values.vars.primary_ip"
```

Inventory-level variables and group resources emitted by the provider are
automatically detected and included in the generated inventory.
