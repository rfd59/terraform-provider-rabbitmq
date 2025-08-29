# Read the exchange settings
data "rabbitmq_exchange_delayed_message" "example" {
  name  = "myexchange"
  vhost = "myvhost"
}
