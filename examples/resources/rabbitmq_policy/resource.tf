# Create a vhost
resource "rabbitmq_vhost" "example" {
  name = "myvhost"
}

# Create a user
resource "rabbitmq_user" "example" {
  name     = "myuser"
  password = "foobar"
  tags     = ["administrator", "management"]
}

# Create a permission
resource "rabbitmq_permissions" "example" {
  user  = rabbitmq_user.example.name
  vhost = rabbitmq_vhost.example.name

  permissions {
    configure = ".*"
    write     = ".*"
    read      = ".*"
  }
}

# Create a policy
resource "rabbitmq_policy" "example" {
  name  = "mypolicy"
  vhost = rabbitmq_permissions.example.vhost

  policy {
    pattern  = ".*"
    priority = 0
    apply_to = "all"

    definition = {
      ha-mode = "all"
    }
  }
}
