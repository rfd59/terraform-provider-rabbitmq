# Create a vhost
resource "rabbitmq_vhost" "example" {
  name = "myvhost"
}

# Create a operator policyst
resource "rabbitmq_operator_policy" "example" {
  name  = "myoperatorpolicy"
  vhost = rabbitmq_vhost.example.vhost

  policy {
    pattern  = ".*"
    priority = 0
    apply_to = "queues"

    definition = {
      message-ttl = 3600000
      expires = 1800000
    }
  }
}
