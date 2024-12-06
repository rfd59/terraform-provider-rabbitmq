<a href="https://terraform.io">
    <img src=".github/tf.png" alt="Terraform logo" title="Terraform" align="left" height="50" />
</a>

# Terraform Provider for RabbitMQ

![GitHub Release](https://img.shields.io/github/v/release/rfd59/terraform-provider-rabbitmq)
![GitHub Issues](https://img.shields.io/github/issues/rfd59/terraform-provider-rabbitmq)
![GitHub Pull Requests](https://img.shields.io/github/issues-pr/rfd59/terraform-provider-rabbitmq)
![GitHub License](https://img.shields.io/github/license/rfd59/terraform-provider-rabbitmq)

![Go version](https://img.shields.io/github/go-mod/go-version/rfd59/terraform-provider-rabbitmq)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/rfd59/terraform-provider-rabbitmq/.github%2Fworkflows%2Fbuild.yml)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/rfd59/terraform-provider-rabbitmq/.github%2Fworkflows%2Ftest.yml?label=test)
![Coverage](https://sonar.rfd.ovh/api/project_badges/measure?project=rfd59.terraform-provider-rabbitmq&metric=coverage&token=sqb_44b6ae8e30de40b0d76cc3bcfad1a5e2e3f3c0c0)

[RabbitMQ](https://rabbitmq.com) is an AMQP message broker server. The **RabbitMQ provider** exposes resources used to manage the configuration of resources in a RabbitMQ server.

The provider supports versions `4.0.x`, `3.13.x` and `3.12.x` of RabbitMQ. It may still work with versions `3.11.x`, `3.10.x`, `3.9.x` and `3.8.x`, however these versions are no longer officialy supported.
> For information on RabbitMQ versions, see the RabbitMQ [Release Information](https://www.rabbitmq.com/release-information).

## Usage Example

```hcl
# 1. Specify the version of the RabbitMQ Provider to use
terraform {
  required_providers {
    rabbitmq = {
      source = "rfd59/rabbitmq"
      version = "2.3.0"
    }
  }
}

# 2. Configure the RabbitMQ Provider
provider "rabbitmq" {
  # The RabbitMQ management plugin must be enabled on the server, to use this provider.
  # You can enable the plugin by doing something similar to `$ sudo rabbitmq-plugins enable rabbitmq_management`
  # https://www.rabbitmq.com/docs/management

  endpoint = "http://127.0.0.1:15672"
  username = "guest"
  password = "guest"
}

# 3. Create a Virtual Hosts into the RabbitMQ Server
resource "rabbitmq_vhost" "example" {
  name = "vhost_example"
}
```

## Developing & Contributing to the Provider

### Requirements

- [Terraform](https://www.terraform.io/downloads.html) 1.0.x
- [Go](https://golang.org/doc/install) 1.23
- [Docker Engine](https://docs.docker.com/engine/install) >= 27.2.x
- [Docker Compose plugin](https://docs.docker.com/compose/install/#scenario-two-install-the-compose-plugin) >= 2.29.x

### Building the Provider

1. Fork and Clone this repository localy
2. To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

   ```sh
   \$ make build
   ...
   \$ $GOPATH/bin/terraform-provider-rabbitmq
   ...
   ```

   > To find the '${GOPATH}', you can run `go env GOPATH`.

3. In order to run the full suite of Acceptance tests, run `make testacc`.

   ```sh
   \$ make testacc
   ```

4. In order to build the Provider documentation, run `make doc`.

   ```sh
   \$ make doc
   ```

### Using this Provider locally

To launch the **examples** terraform scripts with your local Provider, follow these steps:

1. Into your home folder, add/update the `~/.terraformrc` file with:

   ```txt
   provider_installation {

   dev_overrides {
         "rfd59/rabbitmq" = "${GOPATH}/bin"
   }

   # For all other providers, install them directly from their origin provider
   # registries as normal. If you omit this, Terraform will _only_ use
   # the dev_overrides block, and so no other providers will be available.
   direct {}
   }
   ```

   > To find the '${GOPATH}', you can run `go env GOPATH`.

2. Launch a RabbitMQ engine by `./scripts/testacc.sh setup` command.
   > The RabbitMQ console will be available by _http://localhost:15672_ (_guest_/_guest_)
3. If it's not already done, build you local Provider by `make build` command.
4. Launch `terraform -chdir=./examples/... apply` to apply the example.
