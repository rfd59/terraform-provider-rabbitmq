package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
)

type ExchangeResource struct {
	Name     string
	Vhost    string
	Settings ExchangeSettings
}

type ExchangeSettings struct {
	Type              string
	Durable           bool
	AutoDelete        bool
	Internal          bool
	AlternateExchange string
	Arguments         map[string]interface{}
}

func (e *ExchangeResource) RequiredCreate(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"

		settings {}
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

		settings {
			type = "%s"
			durable = %t
			auto_delete = %t
			internal = %t
			alternate_exchange = "%s"
		}
	}
	
	resource "rabbitmq_vhost" "test" {
		name = "%s"
	}
	
	`, data.ResourceType, data.ResourceLabel, e.Name, e.Settings.Type, e.Settings.Durable, e.Settings.AutoDelete, e.Settings.Internal, e.Settings.AlternateExchange, e.Vhost)
}

func (e *ExchangeResource) OptionalUpdate(data TestData) string {
	e.Settings.Type = "fanout"
	e.Settings.AlternateExchange = data.RandomString()

	return e.OptionalCreate(data)
}

func (e *ExchangeResource) OptionalArguments(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"

		settings {
			arguments = {
				"key1" = "%s"
			}
		}
	}`, data.ResourceType, data.ResourceLabel, e.Name, e.Settings.Arguments["key1"])
}

func (e *ExchangeResource) ErrorSettingsBlockStep1(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
	}`, data.ResourceType, data.ResourceLabel, e.Name)
}

func (e *ExchangeResource) ErrorSettingsBlockStep2(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"

		settings {}
		
		settings {
			durable = true
		}
	}`, data.ResourceType, data.ResourceLabel, e.Name)
}

func (e *ExchangeResource) ErrorVhostNotExist(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
		vhost = "%s"

		settings {}
	}`, data.ResourceType, data.ResourceLabel, e.Name, e.Vhost)
}

func (e *ExchangeResource) ErrorAlredayExist(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"

		settings {}
	}
	
	resource "%s" "%s" {
		name = "%s"

		settings {}
	}`, data.ResourceType, data.ResourceLabel, e.Name, data.ResourceType, "same", e.Name)
}

func (q *ExchangeResource) DataSource(data TestData) string {
	return fmt.Sprintf(`
	data "%s" "%s" {
		name = "%s"
	}`, data.ResourceType, data.ResourceLabel, q.Name)
}

func (e ExchangeResource) ExistsInRabbitMQ() error {

	rmqc := TestAcc.Provider.Meta().(*rabbithole.Client)
	myExchange, err := rmqc.GetExchange(e.Vhost, e.Name)
	if err != nil {
		return fmt.Errorf("error retrieving exchange '%s': %#v", e.Name, err)
	}
	if myExchange.Name != e.Name {
		return fmt.Errorf("exchange name is not equal. Actual: '%s' Expected: %s", myExchange.Name, e.Name)
	}
	if myExchange.Vhost != e.Vhost {
		return fmt.Errorf("exchange vhost is not equal. Actual: '%s' Expected: %s", myExchange.Vhost, e.Vhost)
	}
	if myExchange.Type != e.Settings.Type {
		return fmt.Errorf("exchange type is not equal. Actual: '%s' Expected: %s", myExchange.Type, e.Settings.Type)
	}
	if myExchange.Durable != e.Settings.Durable {
		return fmt.Errorf("exchange durable is not equal. Actual: '%t' Expected: %t", myExchange.Durable, e.Settings.Durable)
	}
	if myExchange.AutoDelete != e.Settings.AutoDelete {
		return fmt.Errorf("exchange auto_delete is not equal. Actual: '%t' Expected: %t", myExchange.AutoDelete, e.Settings.AutoDelete)
	}
	if myExchange.Internal != e.Settings.Internal {
		return fmt.Errorf("exchange internal is not equal. Actual: '%t' Expected: %t", myExchange.Internal, e.Settings.Internal)
	}
	lenArg := len(myExchange.Arguments)
	if lenArg > 0 {
		if myExchange.Arguments["alternate-exchange"] == nil {
			if e.Settings.AlternateExchange != "" {
				return fmt.Errorf("exchange alternate_exchange is not equal. Actual: '' Expected: %s", e.Settings.AlternateExchange)
			}
		} else {
			lenArg--
			if myExchange.Arguments["alternate-exchange"] != e.Settings.AlternateExchange {
				return fmt.Errorf("exchange alternate_exchange is not equal. Actual: '%s' Expected: %s", myExchange.Arguments["alternate-exchange"], e.Settings.AlternateExchange)
			}
		}

		if lenArg != len(e.Settings.Arguments) {
			return fmt.Errorf("exchange arguments size is not equal. Actual: '%d' Expected: %d", lenArg, len(e.Settings.Arguments))
		}

		for key, val := range e.Settings.Arguments {
			if myExchange.Arguments[key] != val {
				return fmt.Errorf("exchange argument %q is not equal. Actual: '%s' Expected: %s", key, myExchange.Arguments[key], val)
			}
		}
	}

	return nil
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
		Type:       e.Settings.Type,
		Durable:    e.Settings.Durable,
		AutoDelete: e.Settings.AutoDelete,
		Internal:   e.Settings.Internal,
		Arguments:  e.Settings.Arguments,
	}

	rmqc := TestAcc.Client(t)

	resp, err := rmqc.DeclareExchange(e.Vhost, e.Name, settings)
	if err != nil || resp.StatusCode >= 400 {
		t.Errorf("Failed to init the test!")
	}
}

func (e *ExchangeResource) DelDataSourceExchange(t *testing.T) {
	rmqc := TestAcc.Client(t)

	resp, err := rmqc.DeleteExchange(e.Vhost, e.Name)
	if err != nil || resp.StatusCode >= 400 {
		t.Errorf("Failed to reset the test!")
	}
}
