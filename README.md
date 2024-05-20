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

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 1.0.x
- [Go](https://golang.org/doc/install) 1.21 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-rabbitmq`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-rabbitmq
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-rabbitmq
$ make build
```

Using the provider
------------------

The provider supports versions `3.13.x`, `3.12.x` and `3.11.x` of RabbitMQ. It may still work with versions `3.10.x`, `3.9.x` and `3.8.x`, however these versions are no longer officialy supported.

For information on RabbitMQ versions, see the RabbitMQ [Release Information](https://www.rabbitmq.com/release-information).

Developing the Provider
-----------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.21+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-rabbitmq
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

```sh
$ make testacc
```

To launch the **examples** terraform script with your local Provider, follow these steps:
1. Into your home folder, add/update the `~/.terraformrc` file with:
```
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
3. Build you local Provider by `make build` command.
4. Launch `terraform -chdir=./examples/xxx apply` to apply the example.
   > _xxx_ is the subfolder name from _./examples_. (`terraform -chdir=./examples/vhost apply`)