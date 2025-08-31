package acceptance_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type ExchangeTopicResource struct {
	ExchangeResource
}

func (e ExchangeTopicResource) ExistsInRabbitMQ() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if _, err := e.ExchangeResource.ExistsInRabbitMQ(true); err != nil {
			return err
		} else {
			return nil
		}
	}
}
