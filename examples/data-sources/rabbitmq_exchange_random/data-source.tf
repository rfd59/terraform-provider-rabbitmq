# Read the exchange settings
data "rabbitmq_exchange_random" "example" {
  name  = "myexchange"
  vhost = "myvhost"
}
