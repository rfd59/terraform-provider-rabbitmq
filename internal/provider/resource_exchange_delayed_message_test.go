package provider_test

import (
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
