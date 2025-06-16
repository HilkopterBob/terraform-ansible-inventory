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

<!-- TODO: document new filtering flags and extended provider support when
available -->
