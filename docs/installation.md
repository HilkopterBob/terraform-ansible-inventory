# Installation

1. Clone the repository:

```bash
git clone https://github.com/HilkopterBob/terraform-ansible-inventory.git
cd terraform-ansible-inventory
```

2. Build the binary:

```bash
go build -o terraform-ansible-inventory ./main.go
```

3. (Optional) Install globally:

```bash
go install github.com/HilkopterBob/terraform-ansible-inventory@latest
```

The project has no external dependencies beyond Go itself. Simply ensure a
recent Go toolchain is installed before building.
