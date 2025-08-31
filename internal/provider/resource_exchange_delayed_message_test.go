package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	acceptance_test "github.com/rfd59/terraform-provider-rabbitmq/test/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExchangeDelayedMessage_Required(t *testing.T) {
	data := acceptance_test.BuildTestData("rabbitmq_exchange_delayed_message", "test")
	r := acceptance_test.ExchangeDelayedMessageResource{
		ExchangeResource: acceptance_test.ExchangeResource{
			Name:    data.RandomString(),
			Vhost:   "/",
			Type:    "x-delayed-message",
			Durable: true},
		DelayedType: "direct"}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance_test.TestAcc.PreCheck(t) },
		Providers:    acceptance_test.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config: r.RequiredCreate(data),
				Check: resource.ComposeTestCheckFunc(
					acceptance_test.That(data.ResourceName).Exists(),
					acceptance_test.That(data.ResourceName).Key("id").IsNotEmpty(),
					acceptance_test.That(data.ResourceName).Key("id").HasValue(r.Name+"@"+r.Vhost),
					acceptance_test.That(data.ResourceName).Key("name").HasValue(r.Name),
					acceptance_test.That(data.ResourceName).Key("vhost").HasValue(r.Vhost),
					acceptance_test.That(data.ResourceName).Key("type").HasValue(r.Type),
					acceptance_test.That(data.ResourceName).Key("delayed_type").HasValue(r.DelayedType),
					acceptance_test.That(data.ResourceName).Key("durable").IsBool(r.Durable),
					acceptance_test.That(data.ResourceName).Key("auto_delete").IsBool(r.AutoDelete),
					acceptance_test.That(data.ResourceName).Key("internal").IsBool(r.Internal),
					acceptance_test.That(data.ResourceName).Key("alternate_exchange").DoesNotExist(),
					acceptance_test.That(data.ResourceName).Key("argument.#").DoesNotExist(),
					r.ExistsInRabbitMQ(),
				),
			},
			{
				Config: r.RequiredUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					acceptance_test.That(data.ResourceName).Exists(),
					acceptance_test.That(data.ResourceName).Key("id").IsNotEmpty(),
					acceptance_test.That(data.ResourceName).Key("id").HasValue(r.Name+"@"+r.Vhost),
					acceptance_test.That(data.ResourceName).Key("name").HasValue(r.Name),
					acceptance_test.That(data.ResourceName).Key("vhost").HasValue(r.Vhost),
					acceptance_test.That(data.ResourceName).Key("type").HasValue(r.Type),
					acceptance_test.That(data.ResourceName).Key("delayed_type").HasValue(r.DelayedType),
					acceptance_test.That(data.ResourceName).Key("durable").IsBool(r.Durable),
					acceptance_test.That(data.ResourceName).Key("auto_delete").IsBool(r.AutoDelete),
					acceptance_test.That(data.ResourceName).Key("internal").IsBool(r.Internal),
					acceptance_test.That(data.ResourceName).Key("alternate_exchange").DoesNotExist(),
					acceptance_test.That(data.ResourceName).Key("argument.#").DoesNotExist(),
					r.ExistsInRabbitMQ(),
				),
			},
		},
	})
}

func TestAccExchangeDelayedMessage_Optional(t *testing.T) {
	data := acceptance_test.BuildTestData("rabbitmq_exchange_delayed_message", "test")
	r := acceptance_test.ExchangeDelayedMessageResource{
		ExchangeResource: acceptance_test.ExchangeResource{
			Name:              data.RandomString(),
			Vhost:             data.RandomString(),
			Type:              "x-delayed-message",
			Durable:           false,
			AutoDelete:        true,
			Internal:          true,
			AlternateExchange: data.RandomString()},
		DelayedType: "topic"}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance_test.TestAcc.PreCheck(t) },
		Providers:    acceptance_test.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config: r.OptionalCreate(data),
				Check: resource.ComposeTestCheckFunc(
					acceptance_test.That(data.ResourceName).Exists(),
					acceptance_test.That(data.ResourceName).Key("id").IsNotEmpty(),
					acceptance_test.That(data.ResourceName).Key("id").HasValue(r.Name+"@"+r.Vhost),
					acceptance_test.That(data.ResourceName).Key("name").HasValue(r.Name),
					acceptance_test.That(data.ResourceName).Key("vhost").HasValue(r.Vhost),
					acceptance_test.That(data.ResourceName).Key("type").HasValue(r.Type),
					acceptance_test.That(data.ResourceName).Key("delayed_type").HasValue(r.DelayedType),
					acceptance_test.That(data.ResourceName).Key("durable").IsBool(r.Durable),
					acceptance_test.That(data.ResourceName).Key("auto_delete").IsBool(r.AutoDelete),
					acceptance_test.That(data.ResourceName).Key("internal").IsBool(r.Internal),
					acceptance_test.That(data.ResourceName).Key("alternate_exchange").HasValue(r.AlternateExchange),
					acceptance_test.That(data.ResourceName).Key("argument.#").DoesNotExist(),
					r.ExistsInRabbitMQ(),
				),
			},
			{
				Config: r.OptionalUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					acceptance_test.That(data.ResourceName).Exists(),
					acceptance_test.That(data.ResourceName).Key("id").IsNotEmpty(),
					acceptance_test.That(data.ResourceName).Key("id").HasValue(r.Name+"@"+r.Vhost),
					acceptance_test.That(data.ResourceName).Key("name").HasValue(r.Name),
					acceptance_test.That(data.ResourceName).Key("vhost").HasValue(r.Vhost),
					acceptance_test.That(data.ResourceName).Key("type").HasValue(r.Type),
					acceptance_test.That(data.ResourceName).Key("delayed_type").HasValue(r.DelayedType),
					acceptance_test.That(data.ResourceName).Key("durable").IsBool(r.Durable),
					acceptance_test.That(data.ResourceName).Key("auto_delete").IsBool(r.AutoDelete),
					acceptance_test.That(data.ResourceName).Key("internal").IsBool(r.Internal),
					acceptance_test.That(data.ResourceName).Key("alternate_exchange").HasValue(r.AlternateExchange),
					acceptance_test.That(data.ResourceName).Key("argument.#").DoesNotExist(),
					r.ExistsInRabbitMQ(),
				),
			},
		},
	})
}

func TestAccExchangeDelayedMessage_ArgumentsString(t *testing.T) {
	data := acceptance_test.BuildTestData("rabbitmq_exchange_delayed_message", "test")
	r := acceptance_test.ExchangeDelayedMessageResource{
		ExchangeResource: acceptance_test.ExchangeResource{
			Name:    data.RandomString(),
			Vhost:   "/",
			Type:    "x-delayed-message",
			Durable: true,
			Arguments: []map[string]interface{}{
				{"key": data.RandomString(), "value": data.RandomString(), "type": "string"},
			}},
		DelayedType: "direct"}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance_test.TestAcc.PreCheck(t) },
		Providers:    acceptance_test.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config: r.OptionalArgumentsString(data),
				Check: resource.ComposeTestCheckFunc(
					acceptance_test.That(data.ResourceName).Exists(),
					acceptance_test.That(data.ResourceName).Key("id").IsNotEmpty(),
					acceptance_test.That(data.ResourceName).Key("id").HasValue(r.Name+"@"+r.Vhost),
					acceptance_test.That(data.ResourceName).Key("name").HasValue(r.Name),
					acceptance_test.That(data.ResourceName).Key("vhost").HasValue(r.Vhost),
					acceptance_test.That(data.ResourceName).Key("type").HasValue(r.Type),
					acceptance_test.That(data.ResourceName).Key("delayed_type").HasValue(r.DelayedType),
					acceptance_test.That(data.ResourceName).Key("durable").IsBool(r.Durable),
					acceptance_test.That(data.ResourceName).Key("auto_delete").IsBool(r.AutoDelete),
					acceptance_test.That(data.ResourceName).Key("internal").IsBool(r.Internal),
					acceptance_test.That(data.ResourceName).Key("alternate_exchange").DoesNotExist(),
					acceptance_test.That(data.ResourceName).Key("argument.#").Exists(),
					acceptance_test.That(data.ResourceName).Key("argument.0.key").HasValue(r.Arguments[0]["key"].(string)),
					acceptance_test.That(data.ResourceName).Key("argument.0.value").HasValue(r.Arguments[0]["value"].(string)),
					acceptance_test.That(data.ResourceName).Key("argument.0.type").HasValue(r.Arguments[0]["type"].(string)),
					r.ExistsInRabbitMQ(),
				),
			},
		},
	})
}

func TestAccExchangeDelayedMessage_ArgumentsNumeric(t *testing.T) {
	data := acceptance_test.BuildTestData("rabbitmq_exchange_delayed_message", "test")
	r := acceptance_test.ExchangeDelayedMessageResource{
		ExchangeResource: acceptance_test.ExchangeResource{
			Name:    data.RandomString(),
			Vhost:   "/",
			Type:    "x-delayed-message",
			Durable: true,
			Arguments: []map[string]interface{}{
				{"key": data.RandomString(), "value": data.RandomInteger(), "type": "numeric"},
			}},
		DelayedType: "direct"}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance_test.TestAcc.PreCheck(t) },
		Providers:    acceptance_test.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config: r.OptionalArgumentsNumeric(data),
				Check: resource.ComposeTestCheckFunc(
					acceptance_test.That(data.ResourceName).Exists(),
					acceptance_test.That(data.ResourceName).Key("id").IsNotEmpty(),
					acceptance_test.That(data.ResourceName).Key("id").HasValue(r.Name+"@"+r.Vhost),
					acceptance_test.That(data.ResourceName).Key("name").HasValue(r.Name),
					acceptance_test.That(data.ResourceName).Key("vhost").HasValue(r.Vhost),
					acceptance_test.That(data.ResourceName).Key("type").HasValue(r.Type),
					acceptance_test.That(data.ResourceName).Key("delayed_type").HasValue(r.DelayedType),
					acceptance_test.That(data.ResourceName).Key("durable").IsBool(r.Durable),
					acceptance_test.That(data.ResourceName).Key("auto_delete").IsBool(r.AutoDelete),
					acceptance_test.That(data.ResourceName).Key("internal").IsBool(r.Internal),
					acceptance_test.That(data.ResourceName).Key("alternate_exchange").DoesNotExist(),
					acceptance_test.That(data.ResourceName).Key("argument.#").Exists(),
					acceptance_test.That(data.ResourceName).Key("argument.0.key").HasValue(r.Arguments[0]["key"].(string)),
					acceptance_test.That(data.ResourceName).Key("argument.0.value").HasValue(fmt.Sprintf("%d", r.Arguments[0]["value"].(int))),
					acceptance_test.That(data.ResourceName).Key("argument.0.type").HasValue(r.Arguments[0]["type"].(string)),
					r.ExistsInRabbitMQ(),
				),
			},
		},
	})
}

func TestAccExchangeDelayedMessage_ArgumentsBoolean(t *testing.T) {
	data := acceptance_test.BuildTestData("rabbitmq_exchange_delayed_message", "test")
	r := acceptance_test.ExchangeDelayedMessageResource{
		ExchangeResource: acceptance_test.ExchangeResource{
			Name:    data.RandomString(),
			Vhost:   "/",
			Type:    "x-delayed-message",
			Durable: true,
			Arguments: []map[string]interface{}{
				{"key": data.RandomString(), "value": true, "type": "boolean"},
			}},
		DelayedType: "direct"}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance_test.TestAcc.PreCheck(t) },
		Providers:    acceptance_test.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config: r.OptionalArgumentsBoolean(data),
				Check: resource.ComposeTestCheckFunc(
					acceptance_test.That(data.ResourceName).Exists(),
					acceptance_test.That(data.ResourceName).Key("id").IsNotEmpty(),
					acceptance_test.That(data.ResourceName).Key("id").HasValue(r.Name+"@"+r.Vhost),
					acceptance_test.That(data.ResourceName).Key("name").HasValue(r.Name),
					acceptance_test.That(data.ResourceName).Key("vhost").HasValue(r.Vhost),
					acceptance_test.That(data.ResourceName).Key("type").HasValue(r.Type),
					acceptance_test.That(data.ResourceName).Key("delayed_type").HasValue(r.DelayedType),
					acceptance_test.That(data.ResourceName).Key("durable").IsBool(r.Durable),
					acceptance_test.That(data.ResourceName).Key("auto_delete").IsBool(r.AutoDelete),
					acceptance_test.That(data.ResourceName).Key("internal").IsBool(r.Internal),
					acceptance_test.That(data.ResourceName).Key("alternate_exchange").DoesNotExist(),
					acceptance_test.That(data.ResourceName).Key("argument.#").Exists(),
					acceptance_test.That(data.ResourceName).Key("argument.0.key").HasValue(r.Arguments[0]["key"].(string)),
					acceptance_test.That(data.ResourceName).Key("argument.0.value").HasValue(fmt.Sprintf("%t", r.Arguments[0]["value"].(bool))),
					acceptance_test.That(data.ResourceName).Key("argument.0.type").HasValue(r.Arguments[0]["type"].(string)),
					r.ExistsInRabbitMQ(),
				),
			},
		},
	})
}

func TestAccExchangeDelayedMessage_ArgumentTypeValidation(t *testing.T) {
	data := acceptance_test.BuildTestData("rabbitmq_exchange_delayed_message", "test")
	r := acceptance_test.ExchangeDelayedMessageResource{
		ExchangeResource: acceptance_test.ExchangeResource{
			Name: data.RandomString(),
			Arguments: []map[string]interface{}{
				{"key": data.RandomString(), "value": data.RandomString(), "type": data.RandomString()},
			}},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance_test.TestAcc.PreCheck(t) },
		Providers:    acceptance_test.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      r.OptionalArgumentsString(data),
				ExpectError: regexp.MustCompile("to be one of"),
			},
		},
	})
}

func TestAccExchangeDelayedMessage_DelayedTypeValidation(t *testing.T) {
	data := acceptance_test.BuildTestData("rabbitmq_exchange_delayed_message", "test")
	r := acceptance_test.ExchangeDelayedMessageResource{
		ExchangeResource: acceptance_test.ExchangeResource{
			Name: data.RandomString(),
		},
		DelayedType: data.RandomString()}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance_test.TestAcc.PreCheck(t) },
		Providers:    acceptance_test.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      r.DelayedTypeValidation(data),
				ExpectError: regexp.MustCompile("to be one of"),
			},
		},
	})
}

func TestAccExchangeDelayedMessage_ErrorVhostNotExist(t *testing.T) {
	data := acceptance_test.BuildTestData("rabbitmq_exchange_delayed_message", "test")
	r := acceptance_test.ExchangeDelayedMessageResource{
		ExchangeResource: acceptance_test.ExchangeResource{
			Name:  data.RandomString(),
			Vhost: data.RandomString(),
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance_test.TestAcc.PreCheck(t) },
		Providers:    acceptance_test.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      r.ErrorVhostNotExist(data),
				ExpectError: regexp.MustCompile("vhost_not_found"),
			},
		},
	})
}

func TestAccExchangeDelayedMessage_AlredayExist(t *testing.T) {
	data := acceptance_test.BuildTestData("rabbitmq_exchange_delayed_message", "test")
	r := acceptance_test.ExchangeDelayedMessageResource{
		ExchangeResource: acceptance_test.ExchangeResource{
			Name:    data.RandomString(),
			Vhost:   "/",
			Type:    "x-delayed-message",
			Durable: true},
		DelayedType: "direct"}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance_test.TestAcc.PreCheck(t) },
		Providers:    acceptance_test.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config: r.RequiredCreate(data),
				Check: resource.ComposeTestCheckFunc(
					acceptance_test.That(data.ResourceName).Exists(),
					r.ExistsInRabbitMQ(),
				),
			},
			{
				Config:      r.ErrorAlredayExist(data),
				ExpectError: regexp.MustCompile("exchange already exists"),
			},
		},
	})
}

func TestAccExchangeDelayedMessage_ImportRequired(t *testing.T) {
	data := acceptance_test.BuildTestData("rabbitmq_exchange_delayed_message", "test")
	r := acceptance_test.ExchangeDelayedMessageResource{
		ExchangeResource: acceptance_test.ExchangeResource{
			Name:    data.RandomString(),
			Vhost:   "/",
			Type:    "x-delayed-message",
			Durable: true},
		DelayedType: "direct"}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance_test.TestAcc.PreCheck(t) },
		Providers:    acceptance_test.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config: r.RequiredCreate(data),
				Check: resource.ComposeTestCheckFunc(
					acceptance_test.That(data.ResourceName).Exists(),
					r.ExistsInRabbitMQ(),
				),
			},
			{
				ResourceName:      data.ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
