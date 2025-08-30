package datasources

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider/infras"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider/utils"
)

func Exchange() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Description: "The name of the exchange.",
			Type:        schema.TypeString,
			Required:    true,
		},

		"vhost": {
			Description: "The vhost to create the resource in. Defaults to `/`.",
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "/",
		},

		"type": {
			Description: "The exchange type.",
			Type:        schema.TypeString,
			Computed:    true,
		},

		"durable": {
			Description: "Whether the exchange survives server restarts.",
			Type:        schema.TypeBool,
			Computed:    true,
		},

		"auto_delete": {
			Description: "If `true`, the exchange will delete itself after at least one queue or exchange has been bound to this one, and then all queues or exchanges have been unbound.",
			Type:        schema.TypeBool,
			Computed:    true,
		},

		"internal": {
			Description: "If `true`, clients cannot publish to this exchange directly. It can only be used with exchange to exchange bindings.",
			Type:        schema.TypeBool,
			Computed:    true,
		},

		"alternate_exchange": {
			Description: "If messages to this exchange cannot otherwise be routed, send them to the alternate exchange named here.",
			Type:        schema.TypeString,
			Computed:    true,
		},

		"argument": {
			Description: "The custom argument of the exchange.",
			Type:        schema.TypeSet,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Description: "The argument key.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"value": {
						Description: "The argument value.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"type": {
						Description: "The value type. Possible values are `string`, `numeric`, `boolean` and `list`.",
						Type:        schema.TypeString,
						Computed:    true,
					},
				},
			},
		},
	}
}

func ReadExchange(d *schema.ResourceData, rmqc infras.IRabbitMQInfra) diag.Diagnostics {
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	vhost := d.Get("vhost").(string)

	exchange, err := rmqc.GetExchange(vhost, name)
	if err != nil {
		return diag.Errorf("exchange '%s@%s' is not found: %#v", name, vhost, err)
	}

	d.Set("name", exchange.Name)
	d.Set("vhost", exchange.Vhost)
	d.Set("type", exchange.Type)
	d.Set("durable", exchange.Durable)
	d.Set("auto_delete", exchange.AutoDelete)
	d.Set("internal", exchange.Internal)

	if len(exchange.Arguments) > 0 {
		if val := exchange.Arguments["alternate-exchange"]; val != nil {
			d.Set("alternate_exchange", val)
			delete(exchange.Arguments, "alternate-exchange")
		}

		var args []interface{}
		for key, value := range exchange.Arguments {
			args = append(args, map[string]interface{}{"key": key, "value": fmt.Sprintf("%v", value), "type": utils.GetArgumentType(value)})
		}
		d.Set("argument", args)
	}

	d.SetId(utils.BuildResourceId(name, vhost))

	return diags
}
