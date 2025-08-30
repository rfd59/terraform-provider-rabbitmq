package acceptance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider/utils"
)

type ExchangeResource struct {
	Name              string
	Vhost             string
	Type              string
	Durable           bool
	AutoDelete        bool
	Internal          bool
	AlternateExchange string
	Arguments         []map[string]interface{}
}

func (e *ExchangeResource) RequiredCreate(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
	}`, data.ResourceType, data.ResourceLabel, e.Name)
}

func (e *ExchangeResource) RequiredUpdate(data TestData) string {
	e.Name = data.RandomString()
	return e.RequiredCreate(data)
}

func (e *ExchangeResource) OptionalCreate(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
		vhost = rabbitmq_vhost.test.name

		durable = %t
		auto_delete = %t
		internal = %t
		alternate_exchange = "%s"
	}
	
	resource "rabbitmq_vhost" "test" {
		name = "%s"
	}
	`, data.ResourceType, data.ResourceLabel, e.Name, e.Durable, e.AutoDelete, e.Internal, e.AlternateExchange, e.Vhost)
}

func (e *ExchangeResource) OptionalUpdate(data TestData) string {
	e.AlternateExchange = data.RandomString()

	return e.OptionalCreate(data)
}

func (e *ExchangeResource) OptionalArgumentsString(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"

		argument {
			key = "%s"
			value = "%s"
		    type = "%s"
		}
	}`, data.ResourceType, data.ResourceLabel, e.Name, e.Arguments[0]["key"], e.Arguments[0]["value"], e.Arguments[0]["type"])
}

func (e *ExchangeResource) OptionalArgumentsNumeric(data TestData) string {
	var val string

	switch e.Arguments[0]["value"].(type) {
	case int:
		val = fmt.Sprintf("%d", e.Arguments[0]["value"])
	case float64:
		val = fmt.Sprintf("%.2f", e.Arguments[0]["value"])
	}

	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"

		argument {
			key = "%s"
			value = "%s"
		    type = "%s"
		}
	}`, data.ResourceType, data.ResourceLabel, e.Name, e.Arguments[0]["key"], val, e.Arguments[0]["type"])
}

func (e *ExchangeResource) OptionalArgumentsBoolean(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"

		argument {
			key = "%s"
			value = "%t"
		    type = "%s"
		}
	}`, data.ResourceType, data.ResourceLabel, e.Name, e.Arguments[0]["key"], e.Arguments[0]["value"], e.Arguments[0]["type"])
}

func (e *ExchangeResource) ErrorVhostNotExist(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
		vhost = "%s"
	}`, data.ResourceType, data.ResourceLabel, e.Name, e.Vhost)
}

func (e *ExchangeResource) ErrorAlredayExist(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
	}
	
	resource "%s" "%s" {
		name = "%s"
	}`, data.ResourceType, data.ResourceLabel, e.Name, data.ResourceType, "same", e.Name)
}

func (e *ExchangeResource) DataSource(data TestData) string {
	return fmt.Sprintf(`
	data "%s" "%s" {
		name = "%s"
	}`, data.ResourceType, data.ResourceLabel, e.Name)
}

func (e ExchangeResource) ExistsInRabbitMQ(argsChecked bool) (*rabbithole.DetailedExchangeInfo, error) {
	rmqc := TestAcc.Provider.Meta().(*rabbithole.Client)
	myExchange, err := rmqc.GetExchange(e.Vhost, e.Name)
	if err != nil {
		return nil, fmt.Errorf("error retrieving exchange '%s': %#v", e.Name, err)
	}
	if myExchange.Name != e.Name {
		return nil, fmt.Errorf("exchange 'name' is not equal: expected: '%s', got '%s'", e.Name, myExchange.Name)
	}
	if myExchange.Vhost != e.Vhost {
		return nil, fmt.Errorf("exchange 'vhost' is not equal: expected: '%s', got '%s'", e.Vhost, myExchange.Vhost)
	}
	if myExchange.Type != e.Type {
		return nil, fmt.Errorf("exchange 'type' is not equal: expected: '%s', got '%s'", e.Type, myExchange.Type)
	}
	if myExchange.Durable != e.Durable {
		return nil, fmt.Errorf("exchange 'durable' is not equal: expected: '%t', got '%t'", e.Durable, myExchange.Durable)
	}
	if myExchange.AutoDelete != e.AutoDelete {
		return nil, fmt.Errorf("exchange 'auto_delete' is not equal: expected: '%t', got '%t'", e.AutoDelete, myExchange.AutoDelete)
	}
	if myExchange.Internal != e.Internal {
		return nil, fmt.Errorf("exchange 'internal' is not equal: expected: '%t', got '%t'", e.Internal, myExchange.Internal)
	}
	lenArg := len(myExchange.Arguments)
	if lenArg > 0 {
		if myExchange.Arguments["alternate-exchange"] == nil {
			if e.AlternateExchange != "" {
				return nil, fmt.Errorf("exchange 'alternate_exchange' is not equal: expected '', got '%s'", e.AlternateExchange)
			}
		} else {
			lenArg--
			if myExchange.Arguments["alternate-exchange"] != e.AlternateExchange {
				return nil, fmt.Errorf("exchange 'alternate_exchange' is not equal: expected: '%s', got '%s'", e.AlternateExchange, myExchange.Arguments["alternate-exchange"])
			}
		}

		if argsChecked {
			if lenArg != len(e.Arguments) {
				return nil, fmt.Errorf("exchange arguments size is not equal: expected '%d', got '%d'", len(e.Arguments), lenArg)
			}

			for _, v := range e.Arguments {
				if myExchange.Arguments[v["key"].(string)] != v["value"] {
					return nil, fmt.Errorf("exchange argument %q is not equal: expected: '%s', got '%s'", v["key"], v["value"], myExchange.Arguments[v["key"].(string)])
				}
			}
		}
	}

	return myExchange, nil
}

func (e *ExchangeResource) CheckDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rmqc := TestAcc.Provider.Meta().(*rabbithole.Client)
		exchange, err := rmqc.GetExchange(e.Vhost, e.Name)
		if err != nil && err.(rabbithole.ErrorResponse).StatusCode != 404 {
			return fmt.Errorf("error retrieving exchange '%s': %#v", e.Name, err)
		}

		if exchange != nil {
			return fmt.Errorf("exchange still exists: %s", e.Name)
		}

		return nil
	}
}

func (e *ExchangeResource) SetDataSourceExchange(t *testing.T) {
	settings := rabbithole.ExchangeSettings{
		Type:       e.Type,
		Durable:    e.Durable,
		AutoDelete: e.AutoDelete,
		Internal:   e.Internal,
		Arguments:  map[string]interface{}{},
	}

	if e.AlternateExchange != "" {
		settings.Arguments["alternate-exchange"] = e.AlternateExchange
	}

	for _, v := range e.Arguments {
		if value, err := utils.GetArgumentValue(v); err == nil {
			settings.Arguments[v["key"].(string)] = value
		}
	}

	rmqc := TestAcc.Client(t)
	resp, err := rmqc.DeclareExchange(e.Vhost, e.Name, settings)
	if err != nil || resp.StatusCode >= 400 {
		t.Errorf("Failed to init the test! [%v]", err)
	}
}

func (e *ExchangeResource) DelDataSourceExchange(t *testing.T) {
	rmqc := TestAcc.Client(t)

	resp, err := rmqc.DeleteExchange(e.Vhost, e.Name)
	if err != nil || resp.StatusCode >= 400 {
		t.Errorf("Failed to reset the test!")
	}
}
