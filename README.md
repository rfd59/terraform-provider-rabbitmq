<a href="https://terraform.io">
  <img src=".github/tf.png" alt="Terraform logo" title="Terraform" align="left" height="50" />
</a>

# Terraform Provider for RabbitMQ

![GitHub Release](https://img.shields.io/github/v/release/rfd59/terraform-provider-rabbitmq)
![GitHub Issues](https://img.shields.io/github/issues/rfd59/terraform-provider-rabbitmq)
![GitHub Pull Requests](https://img.shields.io/github/issues-pr/rfd59/terraform-provider-rabbitmq)
![GitHub License](https://img.shields.io/github/license/rfd59/terraform-provider-rabbitmq)

![Go version](https://img.shields.io/github/go-mod/go-version/rfd59/terraform-provider-rabbitmq)
[![Go Report Card](https://goreportcard.com/badge/github.com/terraform-providers/terraform-provider-rabbitmq)](https://goreportcard.com/report/github.com/terraform-providers/terraform-provider-rabbitmq)
![Build Status](https://img.shields.io/github/actions/workflow/status/rfd59/terraform-provider-rabbitmq/.github%2Fworkflows%2Fbuild.yml?label=build)
![Test Status](https://img.shields.io/github/actions/workflow/status/rfd59/terraform-provider-rabbitmq/.github%2Fworkflows%2Ftest.yml?label=test)
[![codecov](https://codecov.io/gh/rfd59/terraform-provider-rabbitmq/graph/badge.svg?token=6DNXX7N0IM)](https://codecov.io/gh/rfd59/terraform-provider-rabbitmq)

---

The **RabbitMQ Terraform Provider** allows you to declaratively manage [RabbitMQ](https://www.rabbitmq.com) resources such as virtual hosts, users, permissions, and more.  

It supports RabbitMQ versions `4.1.x`, `4.0.x`, `3.13.x`, and `3.12.x`. Older releases (`3.11.x` to `3.8.x`) may still work but are **no longer officially supported**.  
➡️ See the official [RabbitMQ release information](https://www.rabbitmq.com/release-information) for details.

---

## 🚀 Quick Start

```hcl
terraform {
  required_providers {
    rabbitmq = {
      source  = "rfd59/rabbitmq"
      version = "2.5.0"
    }
  }
}

provider "rabbitmq" {
  # The RabbitMQ management plugin must be enabled.
  # Enable with: $ sudo rabbitmq-plugins enable rabbitmq_management
  # Docs: https://www.rabbitmq.com/docs/management

  endpoint = "http://127.0.0.1:15672"
  username = "guest"
  password = "guest"
}

resource "rabbitmq_vhost" "example" {
  name = "vhost_example"
}
```

---

## 🛠 Development & Contribution

### Requirements
- [Terraform](https://www.terraform.io/downloads.html) **1.0+**
- [Go](https://golang.org/doc/install) **1.24**
- [Docker Engine](https://docs.docker.com/engine/install) **>= 27.2.x**
- [Docker Compose plugin](https://docs.docker.com/compose/install/#scenario-two-install-the-compose-plugin) **>= 2.29.x**

### Build the Provider

```sh
# Clone your fork
git clone git@github.com:<your-username>/terraform-provider-rabbitmq.git
cd terraform-provider-rabbitmq

# Compile
make build
$GOPATH/bin/terraform-provider-rabbitmq
```

👉 To check your `GOPATH`:  
```sh
go env GOPATH
```

### Run Acceptance Tests

```sh
make testacc
```

### Build Documentation

```sh
make doc
```

---

## 🧪 Using the Provider Locally

1. Add or update your `~/.terraformrc`:

   ```hcl
   provider_installation {
     dev_overrides {
       "rfd59/rabbitmq" = "${GOPATH}/bin"
     }
     direct {}
   }
   ```

2. Start RabbitMQ locally:

   ```sh
   ./scripts/testacc.sh setup
   ```
   ➡️ Console available at [http://localhost:15672](http://localhost:15672) (user: `guest` / password: `guest`)

3. Build the provider:

   ```sh
   make build
   ```

4. Apply examples:

   ```sh
   terraform -chdir=./examples/... apply
   ```

---

## 🤝 Contributing

Contributions are welcome!  

Please open an [issue](https://github.com/rfd59/terraform-provider-rabbitmq/issues) first to discuss what you would like to change.

Guidelines:
- Format code with gofmt
  > `gofmt -w .`
- Lint with golangci-lint
  > `golangci-lint run`
- Write unit tests for new features
- Open a [Pull Request](https://github.com/rfd59/terraform-provider-rabbitmq/pulls) to **develop** branch.    
---

## 📄 License

This project is licensed under the [MIT License](LICENSE).