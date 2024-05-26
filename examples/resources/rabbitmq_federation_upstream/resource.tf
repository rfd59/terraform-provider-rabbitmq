# Create a vhost
resource "rabbitmq_vhost" "example" {
  name = "myvhost"
}

// Create federation upstream
resource "rabbitmq_federation_upstream" "example" {
  name = "myfederationupstream"
  vhost = rabbitmq_vhost.example.vhost

  definition {
    uri = "amqp://guest:guest@upstream-server-name:5672/%2f"
    prefetch_count = 1000
    reconnect_delay = 5
    ack_mode = "on-confirm"
    trust_user_id = false
    max_hops = 1
  }
}
