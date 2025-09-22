package provider

import (
	"context"

	rabbithole "github.com/michaelklishin/rabbit-hole/v3"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcesBindings() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve list of all bindings in a given vhost.",
		ReadContext: dataSourcesListBindings,
		Schema: map[string]*schema.Schema{
			"vhost": {
				Description: "The vhost to list bindings in.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/",
			},
			"bindings": {
				Description: "List of all bindings in the given vhost.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Description: "The source exchange of the binding.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"destination": {
							Description: "The destination exchange of the binding.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"destination_type": {
							Description: "The type of the destination (exchange or queue).",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"routing_key": {
							Description: "The routing key for the binding.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"properties_key": {
							Description: "The properties key for the binding.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"arguments": {
							Description: "The arguments for the binding.",
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourcesListBindings(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	rmqc := meta.(*rabbithole.Client)

	vhost := d.Get("vhost").(string)
	bindings, err := rmqc.ListBindingsIn(vhost)
	if err != nil {
		return diag.FromErr(err)
	}

	var result []map[string]interface{}
	for _, b := range bindings {
		if b.Source != "" {
			result = append(result, map[string]interface{}{
				"source":           b.Source,
				"destination":      b.Destination,
				"destination_type": b.DestinationType,
				"routing_key":      b.RoutingKey,
				"properties_key":   b.PropertiesKey,
				"arguments":        b.Arguments,
			})
		}
	}
	if err := d.Set("bindings", result); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vhost)
	return diags
}
