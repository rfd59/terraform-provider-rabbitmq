package provider_test

import (
	"regexp"
	"testing"

	"github.com/rfd59/terraform-provider-rabbitmq/internal/acceptance"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/acceptance/check"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccQueue_Required(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_queue", "test")
	r := acceptance.QueueResource{Name: data.RandomString(), Vhost: "/", AutoDelete: false, Durable: false}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config: r.RequiredCreate(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("id").MatchesRegex(regexp.MustCompile(r.Name+"@"+r.Vhost)),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("vhost").HasValue(r.Vhost),
					check.That(data.ResourceName).Key("type").HasValue("classic"),
					check.That(data.ResourceName).Key("settings").Count(1),
					check.That(data.ResourceName).Key("settings.0.auto_delete").IsBool(r.AutoDelete),
					check.That(data.ResourceName).Key("settings.0.durable").IsBool(r.Durable),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
		},
	})
}

func TestAccQueue_Optional(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_queue", "test")
	r := acceptance.QueueResource{Name: data.RandomString(), Vhost: "/", AutoDelete: false, Durable: false}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config: r.RequiredCreate(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("id").MatchesRegex(regexp.MustCompile(r.Name+"@"+r.Vhost)),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("vhost").HasValue(r.Vhost),
					check.That(data.ResourceName).Key("type").HasValue("classic"),
					check.That(data.ResourceName).Key("settings").Count(1),
					check.That(data.ResourceName).Key("settings.0.auto_delete").IsBool(r.AutoDelete),
					check.That(data.ResourceName).Key("settings.0.durable").IsBool(r.Durable),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
			{
				Config: r.OptionalUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("id").MatchesRegex(regexp.MustCompile(r.Name+"@"+r.Vhost)),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("vhost").HasValue(r.Vhost),
					check.That(data.ResourceName).Key("settings.0.auto_delete").IsBool(r.AutoDelete),
					check.That(data.ResourceName).Key("settings.0.durable").IsBool(r.Durable),
					check.That(data.ResourceName).Key("settings.0.arguments").DoesNotExist(),
					check.That(data.ResourceName).Key("settings.0.arguments_json").IsEmpty(),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
			{
				Config: r.OptionalUpdateArgument(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("id").MatchesRegex(regexp.MustCompile(r.Name+"@"+r.Vhost)),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("vhost").HasValue(r.Vhost),
					check.That(data.ResourceName).Key("settings.0.auto_delete").IsBool(r.AutoDelete),
					check.That(data.ResourceName).Key("settings.0.durable").IsBool(r.Durable),
					check.That(data.ResourceName).Key("settings.0.arguments.myKey").HasValue("myValue"),
					check.That(data.ResourceName).Key("settings.0.arguments_json").IsEmpty(),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
			{
				Config: r.OptionalUpdateArgumentJson(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("id").MatchesRegex(regexp.MustCompile(r.Name+"@"+r.Vhost)),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("vhost").HasValue(r.Vhost),
					check.That(data.ResourceName).Key("settings.0.auto_delete").IsBool(r.AutoDelete),
					check.That(data.ResourceName).Key("settings.0.durable").IsBool(r.Durable),
					check.That(data.ResourceName).Key("settings.0.arguments").DoesNotExist(),
					check.That(data.ResourceName).Key("settings.0.arguments_json").HasValue("{\"myKey\":\"myValue\"}"),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
		},
	})
}

func TestAccQueue_XQueueType(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_queue", "test")
	r := acceptance.QueueResource{Name: data.RandomString(), Vhost: "/"}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config: r.XQueueTypeArgument(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("id").MatchesRegex(regexp.MustCompile(r.Name+"@"+r.Vhost)),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("vhost").HasValue(r.Vhost),
					check.That(data.ResourceName).Key("type").HasValue("classic"),
					check.That(data.ResourceName).Key("settings.0.arguments.x-queue-type").HasValue("classic"),
					check.That(data.ResourceName).Key("settings.0.arguments_json").IsEmpty(),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
		},
	})
}

func TestAccQueue_ErrorBothArgumentsType(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_queue", "test")
	r := acceptance.QueueResource{Name: data.RandomString(), Vhost: "/"}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      r.ErrorBothArgumentsType(data),
				ExpectError: regexp.MustCompile("Conflicting configuration arguments"),
			},
		},
	})
}

func TestAccQueue_VhostDefaultQueueType(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_queue", "test")
	r := acceptance.QueueResource{Name: data.RandomString(), Vhost: data.RandomString(), Durable: true}

	// Skip the test for RabbitMQ 3.8 & 3.9 : the default queue type is not implemented
	r.SkipTestVhostDefaultQueueType(t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config: r.VhostDefaultQueueType_Step1(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("id").MatchesRegex(regexp.MustCompile(r.Name+"@"+r.Vhost)),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("vhost").HasValue(r.Vhost),
					check.That(data.ResourceName).Key("type").HasValue("quorum"),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
					r.CheckQueueTypeInRabbitMQ("quorum"),
				),
			},
			{
				Config: r.VhostDefaultQueueType_Step2(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("id").MatchesRegex(regexp.MustCompile(r.Name+"@"+r.Vhost)),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("vhost").HasValue(r.Vhost),
					check.That(data.ResourceName).Key("type").HasValue("stream"),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
					r.CheckQueueTypeInRabbitMQ("stream"),
				),
			},
		},
	})
}

func TestAccQueue_AlredayExist(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_queue", "test")
	r := acceptance.QueueResource{Name: data.RandomString(), Vhost: "/"}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config: r.RequiredCreate(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("id").MatchesRegex(regexp.MustCompile(r.Name+"@"+r.Vhost)),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("vhost").HasValue(r.Vhost),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
			{
				Config:      r.AlredayExist(data),
				ExpectError: regexp.MustCompile("queue already exists"),
			},
		},
	})
}

func TestAccQueue_ImportRequired(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_queue", "test")
	r := acceptance.QueueResource{Name: data.RandomString(), Vhost: "/"}

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
