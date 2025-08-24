package core

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	rabbithole "github.com/michaelklishin/rabbit-hole/v3"

	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider/infras"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider/utils"
)

func ResourceExchange() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
						Required:    true,
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
	}
}

func DatasourceExchange() map[string]*schema.Schema {
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

func CreateExchange(d *schema.ResourceData, rmqc infras.IRabbitMQInfra) error {
	name := d.Get("name").(string)
	vhost := d.Get("vhost").(string)

	// Check if already exists
	_, not_found := rmqc.GetExchange(vhost, name)
	if not_found == nil {
		return fmt.Errorf("error creating RabbitMQ exchange '%s': exchange already exists", name)
	}

	// Build exchange info
	info, err := makeInfoExchange(d)
	if err != nil {
		return fmt.Errorf("error creating RabbitMQ exchange '%s': %v", name, err)
	}

	// Declare the exchange
	resp, err := rmqc.DeclareExchange(vhost, name, info)
	if err != nil || resp.StatusCode >= 400 {
		return utils.FailApiResponse(err, resp, "creating", "exchange")
	}

	//Save the id
	d.SetId(utils.BuildResourceId(name, vhost))

	return nil
}

func ReadExchange(d *schema.ResourceData, rmqc infras.IRabbitMQInfra) error {
	name, vhost, err := utils.ParseResourceId(d.Id())
	if err != nil {
		return err
	}

	exchange, err := rmqc.GetExchange(vhost, name)
	if err != nil {
		return utils.CheckDeletedResource(d, err)
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

	return nil
}

func DeleteExchange(d *schema.ResourceData, rmqc infras.IRabbitMQInfra) error {
	name, vhost, err := utils.ParseResourceId(d.Id())
	if err != nil {
		return err
	}

	resp, err := rmqc.DeleteExchange(vhost, name)
	if err != nil || (resp.StatusCode >= 400 && resp.StatusCode != 404) {
		return utils.FailApiResponse(err, resp, "deleting", "exchange")
	}

	return nil
}

func DatasourceReadExchange(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	rmqc := meta.(*rabbithole.Client)
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

func makeInfoExchange(d *schema.ResourceData) (info rabbithole.ExchangeSettings, err error) {
	info.Type = d.Get("type").(string)
	info.Durable = d.Get("durable").(bool)
	info.AutoDelete = d.Get("auto_delete").(bool)
	info.Internal = d.Get("internal").(bool)

	info.Arguments = make(map[string]interface{})
	if v := d.Get("alternate_exchange").(string); len(v) > 0 {
		info.Arguments["alternate-exchange"] = v
	}

	args := d.Get("argument").(*schema.Set)
	for _, v := range args.List() {
		arg := v.(map[string]interface{})
		if value, err := utils.GetArgumentValue(arg); err != nil {
			return rabbithole.ExchangeSettings{}, err
		} else {
			info.Arguments[arg["key"].(string)] = value
		}
	}

	return
}
