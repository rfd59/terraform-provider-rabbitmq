# Create a user (with administrator and management tags)
resource "rabbitmq_user" "example" {
  name     = "myuser"
  password = "foobar"
  tags     = ["administrator", "management"]
}
