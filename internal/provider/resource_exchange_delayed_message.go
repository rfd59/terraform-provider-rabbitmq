package provider

import (
	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider/core"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceExchangeDelayedMessage() *schema.Resource {
	// Load and customize the resource schema
	mySchema := core.SchemaExchange()
	mySchema["delayed_type"] = &schema.Schema{
		Description:  "The type of delayed exchange. Possible values are `direct`, `fanout`, `headers`, `topic`, `x-random` and `x-consistent-hash`. Defaults to `direct`.",
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "direct",
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice([]string{"direct", "fanout", "headers", "topic", "x-random", "x-consistent-hash"}, true),
	}

	return &schema.Resource{
		Description: "The `rabbitmq_exchange_delayed_message` resource creates and manages an _exchange_ of type 'x-delayed-message'.",
		Create:      CreateExchangeDelayedMessage,
		Read:        ReadExchangeDelayedMessage,
		Delete:      DeleteExchangeDelayedMessage,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: mySchema,
	}
}

func CreateExchangeDelayedMessage(d *schema.ResourceData, meta interface{}) error {
	// Set the exchange type
	d.Set("type", "x-delayed-message")

	// Add specific argument
	args := d.Get("argument").(*schema.Set)
	args.Add(map[string]interface{}{"key": "x-delayed-type", "value": d.Get("delayed_type").(string), "type": "string"})
	d.Set("argument", args)

	return core.CreateExchange(d, meta.(*rabbithole.Client))
}

func ReadExchangeDelayedMessage(d *schema.ResourceData, meta interface{}) error {
	if err := core.ReadExchange(d, meta.(*rabbithole.Client)); err != nil {
		return err
	}

	// Add specific argument
	args := d.Get("argument").(*schema.Set)
	for _, v := range args.List() {
		arg := v.(map[string]interface{})
		if arg["key"].(string) == "x-delayed-type" {
			d.Set("delayed_type", arg["value"].(string))
			args.Remove(arg)
			break
		}
	}
	d.Set("argument", args)

	return nil
}

func DeleteExchangeDelayedMessage(d *schema.ResourceData, meta interface{}) error {
	return core.DeleteExchange(d, meta.(*rabbithole.Client))
}
