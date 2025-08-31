# Read the exchange settings
data "rabbitmq_exchange_topic" "example" {
  name  = "myexchange"
  vhost = "myvhost"
}
