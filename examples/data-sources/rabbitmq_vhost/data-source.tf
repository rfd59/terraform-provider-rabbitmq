# Read the vhost settings
data "rabbitmq_vhost" "example" {
  name                = "myvhost"
}
