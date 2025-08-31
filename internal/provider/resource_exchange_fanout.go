package provider

import (
	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider/core/resources"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceExchangeFanout() *schema.Resource {
	return &schema.Resource{
		Description: "Exchange --- The `rabbitmq_exchange_fanout` resource creates and manages an _exchange_ of type 'fanout'.",
		Create:      CreateExchangeFanout,
		Read:        ReadExchangeFanout,
		Delete:      DeleteExchangeFanout,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resources.Exchange(),
	}
}

func CreateExchangeFanout(d *schema.ResourceData, meta interface{}) error {
	// Set the exchange type
	d.Set("type", "fanout")

	return resources.CreateExchange(d, meta.(*rabbithole.Client))
}

func ReadExchangeFanout(d *schema.ResourceData, meta interface{}) error {
	return resources.ReadExchange(d, meta.(*rabbithole.Client))
}

func DeleteExchangeFanout(d *schema.ResourceData, meta interface{}) error {
	return resources.DeleteExchange(d, meta.(*rabbithole.Client))
}
