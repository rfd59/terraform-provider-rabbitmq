package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider/core"
)

func datasourceExchangeDelayedMessage() *schema.Resource {
	// Load and customize the resource schema
	mySchema := core.DatasourceExchange()
	mySchema["delayed_type"] = &schema.Schema{
		Description: "The type of delayed exchange. Possible values are `direct`, `fanout`, `headers`, `topic`, `x-random` and `x-consistent-hash`.",
		Type:        schema.TypeString,
		Computed:    true,
	}

	return &schema.Resource{
		Description: "Exchange --- Use this data source to access information about an existing _exchange_ of type 'x-delayed-message'.",
		ReadContext: datasourceReadExchangeDelayedMessage,
		Schema:      mySchema,
	}
}

func datasourceReadExchangeDelayedMessage(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diag := core.DatasourceReadExchange(ctx, d, meta)

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

	return diag
}
