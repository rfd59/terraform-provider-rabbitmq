# Add rabbitMQ provider
terraform {
  required_providers {
    rabbitmq = {
      source = "rfd59/rabbitmq"
    }
  }
}

# Configure the RabbitMQ provider
provider "rabbitmq" {
  endpoint = "http://127.0.0.1:15672"
  username = "guest"
  password = "guest"
}

# Create a vhost
resource "rabbitmq_vhost" "test" {
  name = "PermissionsTestVhost"
}

# Create a user
resource "rabbitmq_user" "test" {
  name     = "PermissionsTestUser"
  password = "foobar"
  tags     = ["administrator"]
}

# Set Permissions of User to Vhost
resource "rabbitmq_permissions" "test" {
  user  = "${rabbitmq_user.test.name}"
  vhost = "${rabbitmq_vhost.test.name}"

  permissions {
    configure = ".*"
    write     = ".*"
    read      = ".*"
  }
}