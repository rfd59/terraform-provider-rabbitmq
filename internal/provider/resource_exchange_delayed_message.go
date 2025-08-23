package provider

import (
	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider/core"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceExchangeDelayedMessage() *schema.Resource {
	return &schema.Resource{
		Description: "The `rabbitmq_exchange_delayed_message` resource creates and manages an _exchange_ of type 'x-delayed-message'.",
		Create:      CreateExchangeDelayedMessage,
		Read:        ReadExchangeDelayedMessage,
		Delete:      DeleteExchangeDelayedMessage,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the exchange.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"vhost": {
				Description: "The vhost to create the resource in. Defaults to `/`.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/",
				ForceNew:    true,
			},

			"type": {
				Description: "The exchange type.",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"delayed_type": {
				Description:  "The type of delayed exchange. Possible values are `direct`, `fanout`, `headers`, `topic`, `x-random` and `x-consistent-hash`. Defaults to `direct`.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "direct",
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"direct", "fanout", "headers", "topic", "x-random", "x-consistent-hash"}, true),
			},

			"durable": {
				Description: "Whether the exchange survives server restarts. Defaults to `true`.",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     true,
			},

			"auto_delete": {
				Description: "If `true`, the exchange will delete itself after at least one queue or exchange has been bound to this one, and then all queues or exchanges have been unbound. Defaults to `false`.",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
			},

			"internal": {
				Description: "If `true`, clients cannot publish to this exchange directly. It can only be used with exchange to exchange bindings. Defaults to `false`.",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
			},

			"alternate_exchange": {
				Description: "If messages to this exchange cannot otherwise be routed, send them to the alternate exchange named here.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},

			"argument": {
				Description: "The custom argument of the exchange.",
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Description: "The argument key.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"value": {
							Description: "The argument value.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"type": {
							Description:  "The value type. Possible values are `string`, `numeric`, `boolean` and `list`. Defaults to `string`.",
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "string",
							ValidateFunc: validation.StringInSlice([]string{"string", "numeric", "boolean", "list"}, true),
						},
					},
				},
			},
		},
	}
}

func CreateExchangeDelayedMessage(d *schema.ResourceData, meta interface{}) error {
	d.Set("type", "x-delayed-message")

	// Add specific argument
	args := d.Get("argument").(*schema.Set)
	args.Add(map[string]interface{}{"key": "x-delayed-type", "value": d.Get("delayed_type").(string), "type": "string"})
	d.Set("argument", args)

	return core.CreateExchange(d, meta)
}

func ReadExchangeDelayedMessage(d *schema.ResourceData, meta interface{}) error {
	if err := core.ReadExchange(d, meta); err != nil {
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
	return core.DeleteExchange(d, meta)
}
