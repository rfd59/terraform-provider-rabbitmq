# Create a virtual host
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

# Create a queue
resource "rabbitmq_queue" "example" {
  name  = "myqueue"
  vhost = rabbitmq_permissions.example.vhost

  settings {
    durable     = false
    auto_delete = true
    arguments = {
      "x-queue-type" : "quorum",
    }
  }
}
