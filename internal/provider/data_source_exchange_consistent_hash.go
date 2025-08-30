package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider/core"
)

func datasourceExchangeConsistentHash() *schema.Resource {
	return &schema.Resource{
		Description: "Exchange --- Use this data source to access information about an existing _exchange_ of type 'x-consistent-hash'.",
		ReadContext: datasourceReadExchangeConsistentHash,
		Schema:      core.DatasourceExchange(),
	}
}

func datasourceReadExchangeConsistentHash(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return core.DatasourceReadExchange(ctx, d, meta)
}
