# Read the exchange settings
data "rabbitmq_exchange_consistent_hash" "example" {
  name  = "myexchange"
  vhost = "myvhost"
}
