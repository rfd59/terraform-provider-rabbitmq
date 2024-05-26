# Create a virtual host
resource "rabbitmq_vhost" "example" {
  name = "myvhost"
}
