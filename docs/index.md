---
page_title: "RabbitMQ Provider - terraform-provider-rabbitmq"
subcategory: ""
description: |-
  The RabbitMQ provider exposes resources used to manage the configuration of resources in a RabbitMQ server
---

# RabbitMQ Provider

[RabbitMQ](https://rabbitmq.com) is an AMQP message broker server. The RabbitMQ provider exposes resources used to manage the configuration of resources in a RabbitMQ server.

Use the navigation to the left to read about the available resources.

## Example Usage

```terraform
# Configure the RabbitMQ provider
provider "rabbitmq" {
  endpoint = "http://127.0.0.1:15672"
  username = "guest"
  password = "guest"
 
  headers = {
    "X-Custom-Header" = "CustomValue"
  }
}

# Create a virtual host
resource "rabbitmq_vhost" "example" {
  name = "myvhost"
}
```

## Requirements

The [RabbitMQ management plugin](https://www.rabbitmq.com/docs/management) must be enabled on the server, to use this provider. You can enable the plugin by doing something similar to:
```sh
$ sudo rabbitmq-plugins enable rabbitmq_management
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `endpoint` (String) The HTTP URL of the management plugin on the RabbitMQ server. This can also be sourced from the `RABBITMQ_ENDPOINT` Environment Variable.
- `password` (String) Password for the given user. This can also be sourced from the `RABBITMQ_PASSWORD` Environment Variable.
- `username` (String) Username to use to authenticate with the server. This can also be sourced from the `RABBITMQ_USERNAME` Environment Variable.

### Optional

- `cacert_file` (String) The path to a custom CA / intermediate certificate. This can also be sourced from the `RABBITMQ_CACERT` Environment Variable.
- `clientcert_file` (String) The path to the X.509 client certificate. This can also be sourced from the `RABBITMQ_CLIENTCERT` Environment Variable.
- `clientkey_file` (String) The path to the private key. This can also be sourced from the `RABBITMQ_CLIENTKEY` Environment Variable.
- `headers` (Map of String) Custom headers to include in HTTP requests. This should be a map of header names to values.
- `insecure` (Boolean) Trust self-signed certificates. This can also be sourced from the `RABBITMQ_INSECURE` Environment Variable.
- `proxy` (String) The URL of a proxy through which to send HTTP requests to the RabbitMQ server. This can also be sourced from the `RABBITMQ_PROXY` Environment Variable. If not set, the default `HTTP_PROXY`/`HTTPS_PROXY` will be used instead.