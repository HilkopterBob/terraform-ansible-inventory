# Advanced Configuration

The CLI allows overriding JSON paths for hostnames and IP addresses.

```bash
terraform-ansible-inventory \
  --input state.json \
  --format ansible \
  --host-field "values.custom_hostname" \
  --ip-field "values.vars.primary_ip"
```

Use this when your Terraform state stores host information under custom keys.

The provider exposes resources for groups and inventory level variables. These
are automatically detected and included in the generated inventory.
