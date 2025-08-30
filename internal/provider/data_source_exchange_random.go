package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider/core/datasources"
)

func datasourceExchangeRandom() *schema.Resource {
	return &schema.Resource{
		Description: "Exchange --- Use this data source to access information about an existing _exchange_ of type 'x-random'.",
		ReadContext: datasourceReadExchangeRandom,
		Schema:      datasources.Exchange(),
	}
}

func datasourceReadExchangeRandom(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return datasources.ReadExchange(d, meta.(*rabbithole.Client))
}
