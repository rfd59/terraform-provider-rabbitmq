# Read the exchange settings
data "rabbitmq_exchange_headers" "example" {
  name  = "myexchange"
  vhost = "myvhost"
}
