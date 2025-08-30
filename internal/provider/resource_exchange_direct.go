package provider

import (
	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider/core/resources"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceExchangeDirect() *schema.Resource {
	return &schema.Resource{
		Description: "Exchange --- The `rabbitmq_exchange_direct` resource creates and manages an _exchange_ of type 'direct'.",
		Create:      CreateExchangeDirect,
		Read:        ReadExchangeDirect,
		Delete:      DeleteExchangeDirect,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resources.Exchange(),
	}
}

func CreateExchangeDirect(d *schema.ResourceData, meta interface{}) error {
	// Set the exchange type
	d.Set("type", "direct")

	return resources.CreateExchange(d, meta.(*rabbithole.Client))
}

func ReadExchangeDirect(d *schema.ResourceData, meta interface{}) error {
	return resources.ReadExchange(d, meta.(*rabbithole.Client))
}

func DeleteExchangeDirect(d *schema.ResourceData, meta interface{}) error {
	return resources.DeleteExchange(d, meta.(*rabbithole.Client))
}
