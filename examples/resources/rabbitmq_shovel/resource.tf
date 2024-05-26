# Create a virtual host
resource "rabbitmq_vhost" "example" {
    name = "myvhost"
}

# Create an exchange
resource "rabbitmq_exchange" "example" {
    name = "myexchange"
    vhost = rabbitmq_vhost.example.name
    settings {
        type = "fanout"
        durable = false
        auto_delete = true
    }
}

# Create a queue
resource "rabbitmq_queue" "example" {
	name = "myqueue"
	vhost = rabbitmq_vhost.example.name
	settings {
		durable = false
		auto_delete = true
	}
}

# Create a shovel
resource "rabbitmq_shovel" "example" {
	name = "myshovel"
	vhost = rabbitmq_vhost.example.name
	info {
		source_uri = "amqp:///example"
		source_exchange = rabbitmq_exchange.example.name
		source_exchange_key = "example"
		destination_uri = "amqp:///example"
		destination_queue = rabbitmq_queue.example.name
		destination_queue_arguments = {
			x-queue-type = "classic"
		}
	}
}
