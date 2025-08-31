package provider

import (
	"context"
	"fmt"

	rabbithole "github.com/michaelklishin/rabbit-hole/v3"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcesExchange() *schema.Resource {
	return &schema.Resource{
		Description:        "Exchange --- Use this data source to access information about an existing _exchange_.",
		DeprecationMessage: "Migrate this data source to a dedicated exchange data source. This data source will be removed in the next major version of the provider.",
		ReadContext:        dataSourcesReadExchange,
		Schema: map[string]*schema.Schema{
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
				Description: "The vhost to read the exchange in. Defaults to `/`.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/",
			},
			"settings": {
				Description: "The settings of the exchange.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "The type of exchange. Possible values are `direct`, `fanout`, `headers` and `topic`.",
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

						"arguments": {
							Description: "Additional key/value settings for the exchange.",
							Type:        schema.TypeMap,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourcesReadExchange(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	rmqc := meta.(*rabbithole.Client)

	name := d.Get("name").(string)
	vhost := d.Get("vhost").(string)
	id := fmt.Sprintf("%s@%s", name, vhost)

	exchangeSettings, err := rmqc.GetExchange(vhost, name)
	if err != nil {
		return diag.Errorf("exchange '%s@%s' is not found: %#v", name, vhost, err)
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

	d.SetId(id)

	return diags
}
