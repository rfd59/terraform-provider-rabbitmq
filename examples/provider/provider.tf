# Configure the RabbitMQ provider
provider "rabbitmq" {
  endpoint = "http://127.0.0.1:15672"
  username = "guest"
  password = "guest"
 
  headers = {
    "X-Custom-Header" = "CustomValue"
  }
}

# Create a virtual host
resource "rabbitmq_vhost" "example" {
  name = "myvhost"
}