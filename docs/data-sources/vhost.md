---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "rabbitmq_vhost Data Source - terraform-provider-rabbitmq"
subcategory: ""
description: |-
  Use this data source to access information about an existing vhost.
---

# rabbitmq_vhost (Data Source)

Use this data source to access information about an existing _vhost_.

## Example Usage

```terraform
# Read the vhost settings
data "rabbitmq_vhost" "example" {
  name                = "myvhost"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the vhost.

### Read-Only

- `id` (String) The ID of this resource.
