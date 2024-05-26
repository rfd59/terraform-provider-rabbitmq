# Create a vhost
resource "rabbitmq_vhost" "example" {
  name = "myvhost"
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
