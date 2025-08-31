# Create a vhost
resource "rabbitmq_vhost" "example" {
  name = "myvhost"
}

# Create an exchange
resource "rabbitmq_exchange_delayed_message" "example" {
  name  = "myexchange"
  vhost = rabbitmq_vhost.example.vhost

  delayed_type        = "fanout"
  
  durable     = false
  auto_delete = true
  
  argument {
    key   = "myKey"
    value = "12345"
    type  = "numeric"
  }
}
