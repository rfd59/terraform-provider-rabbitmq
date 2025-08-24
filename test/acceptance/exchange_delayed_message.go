package acceptance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type ExchangeDelayedMessageResource struct {
	ExchangeResource
	DelayedType string
}

func (e *ExchangeDelayedMessageResource) OptionalCreate(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
		vhost = rabbitmq_vhost.test.name

		delayed_type = "%s"

		durable = %t
		auto_delete = %t
		internal = %t
		alternate_exchange = "%s"
	}
	
	resource "rabbitmq_vhost" "test" {
		name = "%s"
	}
	`, data.ResourceType, data.ResourceLabel, e.Name, e.DelayedType, e.Durable, e.AutoDelete, e.Internal, e.AlternateExchange, e.Vhost)
}

func (e *ExchangeDelayedMessageResource) OptionalUpdate(data TestData) string {
	e.AlternateExchange = data.RandomString()
	e.DelayedType = "fanout"

	return e.OptionalCreate(data)
}

func (e *ExchangeDelayedMessageResource) DelayedTypeValidation(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
		
		delayed_type = "%s"
	}
	`, data.ResourceType, data.ResourceLabel, e.Name, e.DelayedType)
}

func (e ExchangeDelayedMessageResource) ExistsInRabbitMQ() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if exchange, err := e.ExchangeResource.ExistsInRabbitMQ(false); err != nil {
			return err
		} else {
			if exchange.Arguments["x-delayed-type"] != e.DelayedType {
				return fmt.Errorf("exchange 'delayed_type' is not equal: expected '%s', got '%s'", e.AlternateExchange, exchange.Arguments["x-delayed-type"])
			}
			return nil
		}
	}
}

func (e *ExchangeDelayedMessageResource) SetDataSourceExchange(t *testing.T) {
	// Add x-delayed-type argument
	arg := map[string]interface{}{
		"key":   "x-delayed-type",
		"value": e.DelayedType,
		"type":  "string",
	}
	e.Arguments = append(e.Arguments, arg)
	// Create the exchange
	e.ExchangeResource.SetDataSourceExchange(t)
}
