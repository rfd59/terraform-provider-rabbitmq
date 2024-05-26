# Read the user settings
data "rabbitmq_user" "example" {
  name                = "myuser"
}

# Display the maximum number this user can open
output "max_channels" {
  value = data.rabbitmq_user.example.max_channels
}