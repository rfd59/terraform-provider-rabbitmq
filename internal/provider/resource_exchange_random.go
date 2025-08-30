package provider

import (
	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider/core"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceExchangeRandom() *schema.Resource {
	return &schema.Resource{
		Description: "Exchange --- The `rabbitmq_exchange_random` resource creates and manages an _exchange_ of type 'x-random'.",
		Create:      CreateExchangeRandom,
		Read:        ReadExchangeRandom,
		Delete:      DeleteExchangeRandom,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: core.ResourceExchange(),
	}
}

func CreateExchangeRandom(d *schema.ResourceData, meta interface{}) error {
	// Set the exchange type
	d.Set("type", "x-random")

	return core.CreateExchange(d, meta.(*rabbithole.Client))
}

func ReadExchangeRandom(d *schema.ResourceData, meta interface{}) error {
	return core.ReadExchange(d, meta.(*rabbithole.Client))
}

func DeleteExchangeRandom(d *schema.ResourceData, meta interface{}) error {
	return core.DeleteExchange(d, meta.(*rabbithole.Client))
}
