package provider

import (
	"context"
	"fmt"
	"time"

	rabbithole "github.com/michaelklishin/rabbit-hole/v3"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcesQueue() *schema.Resource {
	return &schema.Resource{
		Description: "Queue --- Use this data source to access information about an existing _queue_.",
		ReadContext: dataSourcesReadQueue,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Description: "The name of the queue.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"vhost": {
				Description: "The virtual host where is stored the queue. Default to `/`.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/",
			},
			"type": {
				Description: "The type of the queue.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"status": {
				Description: "The status of the queue.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourcesReadQueue(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	rmqc := meta.(*rabbithole.Client)

	name := d.Get("name").(string)
	vhost := d.Get("vhost").(string)

	queue, err := rmqc.GetQueue(vhost, name)
	if err != nil {
		return diag.Errorf("queue '%s@%s' is not found: %#v", name, vhost, err)
	}

	d.Set("name", queue.Name)
	d.Set("vhost", queue.Vhost)
	d.Set("type", queue.Type)

	// If the queue is just created, waitting some seconds to have the status
	i := 0
	for queue.Status == "" && i < 10 {
		time.Sleep(time.Second)
		i++
		queue, _ = rmqc.GetQueue(vhost, name)
	}

	d.Set("status", queue.Status)

	d.SetId(fmt.Sprintf("%s@%s", name, vhost))

	return diags
}
