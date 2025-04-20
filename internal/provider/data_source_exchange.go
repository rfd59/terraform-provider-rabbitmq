package provider

import (
	"context"
	"fmt"
	"log"

	rabbithole "github.com/michaelklishin/rabbit-hole/v3"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcesExchange() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to access information about an existing _exchange_.",
		ReadContext: dataSourcesReadExchange,
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
				Description: "The settings of the exchange. The structure is described below.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "The type of exchange.",
							Type:        schema.TypeString,
							Computed:    true,
						},

						"durable": {
							Description: "Whether the exchange survives server restarts.",
							Type:        schema.TypeBool,
							Computed:    true,
						},

						"auto_delete": {
							Description: "Whether the exchange will self-delete when all queues have finished using it.",
							Type:        schema.TypeBool,
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
		return diag.FromErr(checkDeleted(d, err))
	}

	log.Printf("[DEBUG] RabbitMQ: Exchange retrieved %s: %#v", id, exchangeSettings)

	d.Set("name", exchangeSettings.Name)
	d.Set("vhost", exchangeSettings.Vhost)

	exchange := make([]map[string]interface{}, 1)
	e := make(map[string]interface{})
	e["type"] = exchangeSettings.Type
	e["durable"] = exchangeSettings.Durable
	e["auto_delete"] = exchangeSettings.AutoDelete
	e["arguments"] = exchangeSettings.Arguments
	exchange[0] = e
	d.Set("settings", exchange)

	d.SetId(id)

	return diags
}
