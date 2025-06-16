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

<!-- TODO: document provider-specific options like host groups and inventory
level variables when implemented -->
