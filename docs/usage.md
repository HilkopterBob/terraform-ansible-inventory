# Usage

Run the executable with `--help` to see all options:

```bash
terraform-ansible-inventory --help
```

Example commands:

```bash
# Emit raw JSON of all ansible_host objects
terraform-ansible-inventory -i state.json -f json > hosts.json

# Generate INI style inventory
terraform-ansible-inventory -i state.json -f ini > hosts.ini

# Plain text dump
terraform-ansible-inventory -i state.json -f txt > hosts.txt

# Native Ansible inventory
terraform-ansible-inventory -i state.json -f ansible
```

You can restrict the output to specific hosts or groups using the `--host` and
`--group` flags. Multiple values may be provided.
