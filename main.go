package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/rabbitmq"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: rabbitmq.Provider,
	})
}
