# Read the exchange settings
data "rabbitmq_exchange_direct" "example" {
  name  = "myexchange"
  vhost = "myvhost"
}
