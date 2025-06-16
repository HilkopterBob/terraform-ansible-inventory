# Terraform Ansible Inventory

## About

`terraform-ansible-inventory` is a lightweight CLI that extracts host and group information from Terraform state files produced by the `ansible/ansible` provider. It was created to provide a reliable way to generate dynamic inventories directly from Terraform.

The dynamic inventory plugin shipped by Red Hat lacks support for child modules and complex states, so it often fails on real projects. This tool understands nested modules and produces complete inventories with minimal effort.

## Features

- **Streaming JSON parsing** for huge state files with very low memory consumption
- Handles `yaml`, `ini`, `json`, plain text, and native Ansible inventory output
- Supports custom JSON dot-paths for hostnames and IP addresses
- Automatically strips CIDR suffixes when generating `ansible_host` variables
- Full awareness of all resources emitted by the `ansible/ansible` provider, including host groups and inventory variables
- Recognises group hierarchy and child modules in the state
- Simple CLI interface with sensible defaults and a built-in help message
- Includes smoke test data and CI workflow for reliability

## Installation

A recent Go toolchain is the only requirement.

```bash
# Clone and build from source
git clone https://github.com/HilkopterBob/terraform-ansible-inventory.git
cd terraform-ansible-inventory

go build -o terraform-ansible-inventory ./main.go
```

You can also install globally:

```bash
go install github.com/HilkopterBob/terraform-ansible-inventory@latest
```

Pre-built binaries are available on the [releases page](https://github.com/HilkopterBob/terraform-ansible-inventory/releases) for direct download.

## Usage

Run with `--help` to see all options:

```bash
terraform-ansible-inventory --help
```

Common examples:

```bash
# YAML inventory
terraform-ansible-inventory -i state.json -f yaml > inventory.yml

# INI inventory
terraform-ansible-inventory -i state.json -f ini > inventory.ini

# JSON inventory
terraform-ansible-inventory -i state.json -f json > inventory.json

# Native Ansible format
terraform-ansible-inventory -i state.json -f ansible
```

Inventory can be limited to particular hosts or groups via the `--host` and `--group` flags.

## Advanced Configuration

Customise the JSON paths used for hostnames and IPs:

```bash
terraform-ansible-inventory \
  --input state.json \
  --format ansible \
  --host-field "values.custom_hostname" \
  --ip-field "values.vars.primary_ip"
```

## Important Notes

This tool expects your Terraform state to contain the resources exported by the `ansible/ansible` provider. You must use that provider in your Terraform configuration so the necessary host and group data is written to the state file.

## Motivation

The goal of this project is to make it easy to build accurate Ansible inventories from Terraform. Red Hat's official dynamic inventory plugin cannot parse states containing child modules and therefore misses many resources. `terraform-ansible-inventory` parses the state directly and supports complex module layouts, giving you a complete and reliable inventory every time.

