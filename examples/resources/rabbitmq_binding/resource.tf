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

# Create an exchange
resource "rabbitmq_exchange" "example" {
  name  = "myexchange"
  vhost = rabbitmq_vhost.example.vhost

  settings {
    type        = "fanout"
    durable     = false
    auto_delete = true
  }
}

# Create a queue
resource "rabbitmq_queue" "test" {
  name  = "myqueue"
  vhost = rabbitmq_vhost.example.vhost

  settings {
    durable     = true
    auto_delete = false
  }
}

# Create a binding
resource "rabbitmq_binding" "example" {
  source           = rabbitmq_exchange.example.name
  vhost            = rabbitmq_vhost.example.name
  destination      = rabbitmq_queue.example.name
  destination_type = "queue"
  routing_key      = "#"
}
