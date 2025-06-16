# terraform-ansible-inventory CLI

A high-performance, zero-dependency Go command-line tool for extracting Ansible host definitions from deeply nested Terraform state (or any large JSON), and emitting them in various formats (JSON, INI, TXT, or native Ansible inventory).

---

## 🔍 Features

- **Streaming JSON parsing**: Handles files of 20M+ lines and 500+ levels deep without blowing memory.
- **Flexible output formats**: `json`, `ini`, `txt`, or native `ansible` inventory.
- **Extensible dot-path extraction**: Pull any string field from the `values` subtree via `--host-field` and `--ip-field` flags.
- **Built-in IP/CIDR handling**: Strips CIDR suffix for Ansible inventories.
- **Clean CLI interface**: Automatic `--help` message, sensible defaults.
- **CI/CD ready**: Smoke-test JSON example and GitHub Actions for build & test.

---

## 📦 Installation

```bash
# Clone the repo:
git clone https://github.com/HilkopterBob/terraform-ansible-inventory.git
cd terraform-ansible-inventory

# Build the binary:
go build -o terraform-ansible-inventory ./main.go

# (Optional) Install globally:
go install github.com/HilkopterBob/terraform-ansible-inventory@latest
```

---

## 🚀 Usage

```
$ terraform-ansible-inventory --help
Usage of terraform-ansible-inventory:
  -format string
        Output format: json, ini, txt, ansible (default "json")
  -host-field string
        Dot-path to hostname in each object (default "values.name")
  -input string
        Path to input JSON file (or '-' for stdin)
  -ip-field string
        Dot-path to IP in each object (CIDR will be stripped) (default "values.variables.ip")
```

All flags can be abbreviated:

- `-i` or `--input`
- `-f` or `--format`

### Examples

Assume `state.json` is your Terraform state:

```bash
# 1) Default: emit raw JSON array of all ansible_host objects
terraform-ansible-inventory -i state.json -f json > hosts.json

# 2) Generate an INI-style inventory
terraform-ansible-inventory -i state.json -f ini > hosts.ini

# 3) Plain-text dump
terraform-ansible-inventory -i state.json -f txt > hosts.txt

# 4) Native Ansible inventory
terraform-ansible-inventory -i state.json -f ansible
# → example output:
# ns0.example.com ansible_host=10.0.0.1
# ns1.example.com ansible_host=10.0.0.2
```

### Custom field paths

If your JSON uses different keys, override the defaults:

```bash
terraform-ansible-inventory \
  --input state.json \
  --format ansible \
  --host-field "values.custom_hostname" \
  --ip-field   "values.vars.primary_ip"
```

---

## 🔧 Contributing

1. Fork & clone the repo
2. Create a feature branch `git checkout -b feature/myfeature`
3. Add tests under `internal/*_test.go`
4. Ensure all tests pass
5. Submit a PR

---

## 📄 License

Apache-2.0 © Nick von Podewils

