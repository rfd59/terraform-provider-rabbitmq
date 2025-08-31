# Read the exchange settings
data "rabbitmq_exchange_fanout" "example" {
  name  = "myexchange"
  vhost = "myvhost"
}
