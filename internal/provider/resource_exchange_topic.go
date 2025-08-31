package provider

import (
	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider/core/resources"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceExchangeTopic() *schema.Resource {
	return &schema.Resource{
		Description: "Exchange --- The `rabbitmq_exchange_topic` resource creates and manages an _exchange_ of type 'topic'.",
		Create:      CreateExchangeTopic,
		Read:        ReadExchangeTopic,
		Delete:      DeleteExchangeTopic,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resources.Exchange(),
	}
}

func CreateExchangeTopic(d *schema.ResourceData, meta interface{}) error {
	// Set the exchange type
	d.Set("type", "topic")

	return resources.CreateExchange(d, meta.(*rabbithole.Client))
}

func ReadExchangeTopic(d *schema.ResourceData, meta interface{}) error {
	return resources.ReadExchange(d, meta.(*rabbithole.Client))
}

func DeleteExchangeTopic(d *schema.ResourceData, meta interface{}) error {
	return resources.DeleteExchange(d, meta.(*rabbithole.Client))
}
