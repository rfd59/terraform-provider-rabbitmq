package provider

import (
	"fmt"

	rabbithole "github.com/michaelklishin/rabbit-hole/v3"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceExchange() *schema.Resource {
	return &schema.Resource{
		Description:        "The `rabbitmq_exchange` resource creates and manages an exchange.",
		DeprecationMessage: "Migrate this resource to a dedicated exchange resource. This resource will be removed in the next major version of the provider.",
		Create:             CreateExchange,
		Read:               ReadExchange,
		Delete:             DeleteExchange,
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

			"settings": {
				Description: "The settings of the exchange.",
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description:  "The type of exchange. Possible values are `direct`, `fanout`, `headers` and `topic`. Defaults to `direct`.",
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "direct",
							ValidateFunc: validation.StringInSlice([]string{"direct", "fanout", "headers", "topic"}, true),
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

						"arguments": {
							Description: "Additional key/value settings for the exchange.",
							Type:        schema.TypeMap,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func CreateExchange(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	name := d.Get("name").(string)
	vhost := d.Get("vhost").(string)

	// Check if already exists
	_, not_found := rmqc.GetExchange(vhost, name)
	if not_found == nil {
		return fmt.Errorf("error creating RabbitMQ exchange '%s': exchange already exists", name)
	}

	settings := d.Get("settings").([]interface{})[0].(map[string]interface{})
	if err := declareExchange(rmqc, vhost, name, settings); err != nil {
		return err
	}

	id := fmt.Sprintf("%s@%s", name, vhost)
	d.SetId(id)

	return ReadExchange(d, meta)
}

func ReadExchange(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	name, vhost, err := parseResourceId(d)
	if err != nil {
		return err
	}

	exchangeSettings, err := rmqc.GetExchange(vhost, name)
	if err != nil {
		return checkDeleted(d, err)
	}

	d.Set("name", exchangeSettings.Name)
	d.Set("vhost", exchangeSettings.Vhost)

	settingsList := make([]map[string]interface{}, 1)

	settings := make(map[string]interface{})
	settings["type"] = exchangeSettings.Type
	settings["durable"] = exchangeSettings.Durable
	settings["auto_delete"] = exchangeSettings.AutoDelete
	settings["internal"] = exchangeSettings.Internal
	settings["alternate_exchange"] = exchangeSettings.Arguments["alternate-exchange"]
	delete(exchangeSettings.Arguments, "alternate-exchange")
	settings["arguments"] = exchangeSettings.Arguments

	settingsList[0] = settings
	d.Set("settings", settingsList)

	return nil
}

func DeleteExchange(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	name, vhost, err := parseResourceId(d)
	if err != nil {
		return err
	}

	resp, err := rmqc.DeleteExchange(vhost, name)
	if err != nil || (resp.StatusCode >= 400 && resp.StatusCode != 404) {
		return failApiResponse(err, resp, "deleting", "exchange")
	}

	return nil
}

func declareExchange(rmqc *rabbithole.Client, vhost string, name string, settings map[string]interface{}) error {
	exchangeSettings := rabbithole.ExchangeSettings{}

	if v, ok := settings["type"].(string); ok {
		exchangeSettings.Type = v
	}

	if v, ok := settings["durable"].(bool); ok {
		exchangeSettings.Durable = v
	}

	if v, ok := settings["auto_delete"].(bool); ok {
		exchangeSettings.AutoDelete = v
	}

	if v, ok := settings["internal"].(bool); ok {
		exchangeSettings.Internal = v
	}

	if v, ok := settings["arguments"].(map[string]interface{}); ok {
		exchangeSettings.Arguments = v
	}

	if v, ok := settings["alternate_exchange"].(string); ok && len(v) > 0 {
		exchangeSettings.Arguments["alternate-exchange"] = v
	}

	resp, err := rmqc.DeclareExchange(vhost, name, exchangeSettings)
	if err != nil || resp.StatusCode >= 400 {
		return failApiResponse(err, resp, "creating", "exchange")
	}

	return nil
}
