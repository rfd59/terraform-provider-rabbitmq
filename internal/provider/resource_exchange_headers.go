package provider

import (
	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider/core/resources"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceExchangeHeaders() *schema.Resource {
	return &schema.Resource{
		Description: "Exchange --- The `rabbitmq_exchange_headers` resource creates and manages an _exchange_ of type 'headers'.",
		Create:      CreateExchangeHeaders,
		Read:        ReadExchangeHeaders,
		Delete:      DeleteExchangeHeaders,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resources.Exchange(),
	}
}

func CreateExchangeHeaders(d *schema.ResourceData, meta interface{}) error {
	// Set the exchange type
	d.Set("type", "headers")

	return resources.CreateExchange(d, meta.(*rabbithole.Client))
}

func ReadExchangeHeaders(d *schema.ResourceData, meta interface{}) error {
	return resources.ReadExchange(d, meta.(*rabbithole.Client))
}

func DeleteExchangeHeaders(d *schema.ResourceData, meta interface{}) error {
	return resources.DeleteExchange(d, meta.(*rabbithole.Client))
}
