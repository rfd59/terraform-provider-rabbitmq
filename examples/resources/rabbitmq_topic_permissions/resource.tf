# Create a virtual host
resource "rabbitmq_vhost" "example" {
  name = "myvhost"
}

# Create a user
resource "rabbitmq_user" "example" {
  name     = "myuser"
  password = "foobar"
  tags     = ["administrator"]
}

# Create a topic permission
resource "rabbitmq_topic_permissions" "example" {
  user  = rabbitmq_user.example.name
  vhost = rabbitmq_vhost.example.name

  permissions {
    exchange = "amq.topic"
    write    = ".*"
    read     = ".*"
  }
}
