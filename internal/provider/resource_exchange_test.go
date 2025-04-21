package provider_test

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/rfd59/terraform-provider-rabbitmq/internal/acceptance"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/acceptance/check"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExchange_Required(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_exchange", "test")
	r := acceptance.ExchangeResource{Name: data.RandomString(), Vhost: "/", Settings: acceptance.ExchangeSettings{Type: "direct", Durable: true}}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config: r.RequiredCreate(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("id").IsNotEmpty(),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("vhost").HasValue(r.Vhost),
					check.That(data.ResourceName).Key("id").HasValue(r.Name+"@"+r.Vhost),
					check.That(data.ResourceName).Key("settings").Count(1),
					check.That(data.ResourceName).Key("settings.0.type").HasValue(r.Settings.Type),
					check.That(data.ResourceName).Key("settings.0.durable").HasValue(strconv.FormatBool(r.Settings.Durable)),
					check.That(data.ResourceName).Key("settings.0.auto_delete").HasValue(strconv.FormatBool(r.Settings.AutoDelete)),
					check.That(data.ResourceName).Key("settings.0.internal").HasValue(strconv.FormatBool(r.Settings.Internal)),
					check.That(data.ResourceName).Key("settings.0.alternate_exchange").IsEmpty(),
					check.That(data.ResourceName).Key("settings.0.arguments").DoesNotExist(),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
			{
				Config: r.RequiredUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("id").IsNotEmpty(),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("vhost").HasValue(r.Vhost),
					check.That(data.ResourceName).Key("id").HasValue(r.Name+"@"+r.Vhost),
					check.That(data.ResourceName).Key("settings").Count(1),
					check.That(data.ResourceName).Key("settings.0.type").HasValue(r.Settings.Type),
					check.That(data.ResourceName).Key("settings.0.durable").HasValue(strconv.FormatBool(r.Settings.Durable)),
					check.That(data.ResourceName).Key("settings.0.auto_delete").HasValue(strconv.FormatBool(r.Settings.AutoDelete)),
					check.That(data.ResourceName).Key("settings.0.internal").HasValue(strconv.FormatBool(r.Settings.Internal)),
					check.That(data.ResourceName).Key("settings.0.alternate_exchange").IsEmpty(),
					check.That(data.ResourceName).Key("settings.0.arguments").DoesNotExist(),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
		},
	})
}

func TestAccExchange_Optional(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_exchange", "test")
	r := acceptance.ExchangeResource{Name: data.RandomString(), Vhost: data.RandomString(), Settings: acceptance.ExchangeSettings{Type: "topic", Durable: false, AutoDelete: true, Internal: true, AlternateExchange: data.RandomString()}}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config: r.OptionalCreate(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("id").IsNotEmpty(),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("vhost").HasValue(r.Vhost),
					check.That(data.ResourceName).Key("id").HasValue(r.Name+"@"+r.Vhost),
					check.That(data.ResourceName).Key("settings").Count(1),
					check.That(data.ResourceName).Key("settings.0.type").HasValue(r.Settings.Type),
					check.That(data.ResourceName).Key("settings.0.durable").HasValue(strconv.FormatBool(r.Settings.Durable)),
					check.That(data.ResourceName).Key("settings.0.auto_delete").HasValue(strconv.FormatBool(r.Settings.AutoDelete)),
					check.That(data.ResourceName).Key("settings.0.internal").HasValue(strconv.FormatBool(r.Settings.Internal)),
					check.That(data.ResourceName).Key("settings.0.alternate_exchange").HasValue(r.Settings.AlternateExchange),
					check.That(data.ResourceName).Key("settings.0.arguments").DoesNotExist(),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
			{
				Config: r.OptionalUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("id").IsNotEmpty(),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("vhost").HasValue(r.Vhost),
					check.That(data.ResourceName).Key("id").HasValue(r.Name+"@"+r.Vhost),
					check.That(data.ResourceName).Key("settings").Count(1),
					check.That(data.ResourceName).Key("settings.0.type").HasValue(r.Settings.Type),
					check.That(data.ResourceName).Key("settings.0.durable").HasValue(strconv.FormatBool(r.Settings.Durable)),
					check.That(data.ResourceName).Key("settings.0.auto_delete").HasValue(strconv.FormatBool(r.Settings.AutoDelete)),
					check.That(data.ResourceName).Key("settings.0.internal").HasValue(strconv.FormatBool(r.Settings.Internal)),
					check.That(data.ResourceName).Key("settings.0.alternate_exchange").HasValue(r.Settings.AlternateExchange),
					check.That(data.ResourceName).Key("settings.0.arguments").DoesNotExist(),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
		},
	})
}

func TestAccExchange_Arguments(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_exchange", "test")
	r := acceptance.ExchangeResource{Name: data.RandomString(), Vhost: "/", Settings: acceptance.ExchangeSettings{Type: "direct", Durable: true, Arguments: map[string]interface{}{"key1": data.RandomString()}}}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config: r.OptionalArguments(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("id").IsNotEmpty(),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("vhost").HasValue(r.Vhost),
					check.That(data.ResourceName).Key("id").HasValue(r.Name+"@"+r.Vhost),
					check.That(data.ResourceName).Key("settings").Count(1),
					check.That(data.ResourceName).Key("settings.0.type").HasValue(r.Settings.Type),
					check.That(data.ResourceName).Key("settings.0.durable").HasValue(strconv.FormatBool(r.Settings.Durable)),
					check.That(data.ResourceName).Key("settings.0.auto_delete").HasValue(strconv.FormatBool(r.Settings.AutoDelete)),
					check.That(data.ResourceName).Key("settings.0.internal").HasValue(strconv.FormatBool(r.Settings.Internal)),
					check.That(data.ResourceName).Key("settings.0.alternate_exchange").HasValue(r.Settings.AlternateExchange),
					check.That(data.ResourceName).Key("settings.0.arguments.%").IsNotEmpty(),
					check.That(data.ResourceName).Key("settings.0.arguments.key1").HasValue(r.Settings.Arguments["key1"].(string)),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
		},
	})
}

func TestAccExchange_ErrorSettingsBlock(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_exchange", "test")
	r := acceptance.ExchangeResource{Name: data.RandomString()}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      r.ErrorSettingsBlockStep1(data),
				ExpectError: regexp.MustCompile("Insufficient settings blocks"),
			},
			{
				Config:      r.ErrorSettingsBlockStep2(data),
				ExpectError: regexp.MustCompile("Too many settings blocks"),
			},
		},
	})
}

func TestAccExchange_ErrorVhostNotExist(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_exchange", "test")
	r := acceptance.ExchangeResource{Name: data.RandomString(), Vhost: data.RandomString()}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      r.ErrorVhostNotExist(data),
				ExpectError: regexp.MustCompile("vhost_not_found"),
			},
		},
	})
}

func TestAccExchange_AlredayExist(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_exchange", "test")
	r := acceptance.ExchangeResource{Name: data.RandomString(), Vhost: "/", Settings: acceptance.ExchangeSettings{Type: "direct", Durable: true}}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config: r.RequiredCreate(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
			{
				Config:      r.ErrorAlredayExist(data),
				ExpectError: regexp.MustCompile("exchange already exists"),
			},
		},
	})
}

func TestAccExchange_ImportRequired(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_exchange", "test")
	r := acceptance.ExchangeResource{Name: data.RandomString(), Vhost: "/", Settings: acceptance.ExchangeSettings{Type: "direct", Durable: true}}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config: r.RequiredCreate(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
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
