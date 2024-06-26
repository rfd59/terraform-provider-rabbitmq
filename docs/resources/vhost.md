---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "rabbitmq_vhost Resource - terraform-provider-rabbitmq"
subcategory: ""
description: |-
  The rabbitmq_vhost resource creates and manages a vhost.
---

# rabbitmq_vhost (Resource)

The `rabbitmq_vhost` resource creates and manages a vhost.

## Example Usage

```terraform
# Create a virtual host
resource "rabbitmq_vhost" "example" {
  name = "myvhost"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the vhost.

### Optional

- `default_queue_type` (String) Default queue type for new queues. The available values are `classic`, `quorum` or `stream`.
- `description` (String) A friendly description.
- `max_connections` (String) To limit the total number of concurrent client connections in vhost.
- `max_queues` (String) To limit the total number of queues in vhost.
- `tracing` (Boolean) To enable/disable tracing.

### Read-Only

- `id` (String) The ID of this resource.
