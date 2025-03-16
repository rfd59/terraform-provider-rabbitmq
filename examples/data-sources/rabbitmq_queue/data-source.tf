# Read the queue settings
data "rabbitmq_queue" "example" {
  name                = "myqueue"
}

# Display the queue type
output "type" {
  value = data.rabbitmq_queue.example.type
}
