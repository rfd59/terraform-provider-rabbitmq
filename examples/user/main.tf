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

# Create a user
resource "rabbitmq_user" "test" {
  name     = "UserTest"
  password = "foobar"
  tags     = ["administrator", "management"]
}