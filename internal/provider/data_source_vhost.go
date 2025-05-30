package provider

import (
	"context"

	rabbithole "github.com/michaelklishin/rabbit-hole/v3"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcesVhost() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to access information about an existing _vhost_.",
		ReadContext: dataSourcesReadVhost,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Description: "The name of the vhost.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourcesReadVhost(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	rmqc := meta.(*rabbithole.Client)

	name := d.Get("name").(string)

	vhost, err := rmqc.GetVhost(name)
	if err != nil {
		return diag.Errorf("vhost '%s' is not found: %#v", name, err)
	}

	d.Set("name", vhost.Name)

	d.SetId(name)

	return diags
}
