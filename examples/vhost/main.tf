# Add rabbitMQ provider
terraform {
  required_providers {
    rabbitmq = {
      source = "rfd59/rabbitmq"
    }
  }
}

# Configure the RabbitMQ provider
provider "rabbitmq" {
  endpoint = "http://127.0.0.1:15672"
  username = "guest"
  password = "guest"
}

# Create a basic virtual host
resource "rabbitmq_vhost" "test_basic" {
  name = "VhostTestBasic"
}

# Create a full virtual host
resource "rabbitmq_vhost" "test_full" {
  name = "VhostTestFull"
  description = "A vhost with full settings..."
  default_queue_type = "classic"
  max_connections = 10
  max_queues = 10
}