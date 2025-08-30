package provider

import (
	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider/core"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceExchangeConsistentHash() *schema.Resource {
	return &schema.Resource{
		Description: "Exchange --- The `rabbitmq_exchange_consistent_hash` resource creates and manages an _exchange_ of type 'x-consistent-hash'.",
		Create:      CreateExchangeConsistentHash,
		Read:        ReadExchangeConsistentHash,
		Delete:      DeleteExchangeConsistentHash,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: core.ResourceExchange(),
	}
}

func CreateExchangeConsistentHash(d *schema.ResourceData, meta interface{}) error {
	// Set the exchange type
	d.Set("type", "x-consistent-hash")

	return core.CreateExchange(d, meta.(*rabbithole.Client))
}

func ReadExchangeConsistentHash(d *schema.ResourceData, meta interface{}) error {
	return core.ReadExchange(d, meta.(*rabbithole.Client))
}

func DeleteExchangeConsistentHash(d *schema.ResourceData, meta interface{}) error {
	return core.DeleteExchange(d, meta.(*rabbithole.Client))
}
