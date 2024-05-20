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

# Create a vhost
resource "rabbitmq_vhost" "test" {
  name = "QueueTestVhost"
}

# Create a basic queue
resource "rabbitmq_queue" "test_basic" {
  name  = "QueueTestBasic"
  vhost = "${rabbitmq_vhost.test.name}"

  settings {}
}

# Create a classic queue
resource "rabbitmq_queue" "test_classic" {
  name  = "QueueTestClassic"
  vhost = "${rabbitmq_vhost.test.name}"

  settings {
    durable     = false
    auto_delete = true
    arguments = {
      "x-max-length-bytes" : 45
    }
  }
}

# Create a quorum queue
resource "rabbitmq_queue" "test_quorum" {
  name  = "QueueTestQuorum"
  vhost = "${rabbitmq_vhost.test.name}"

  settings {
    arguments = {
      "x-queue-type" : "quorum",
      "x-overflow": "drop-head"
    }
  }
}

# Create a stream queue
resource "rabbitmq_queue" "test_stream" {
  name  = "QueueTestStream"
  vhost = "${rabbitmq_vhost.test.name}"

  settings {
    arguments = {
      "x-queue-type" : "stream",
      "x-max-age" : "1h"
    }
  }
}
