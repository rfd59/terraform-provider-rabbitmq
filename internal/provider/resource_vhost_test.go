package provider_test

import (
	"regexp"
	"testing"

	"github.com/rfd59/terraform-provider-rabbitmq/internal/acceptance"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/acceptance/check"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVhost_Required(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_vhost", "test")
	r := acceptance.VhostResource{Name: data.RandomString()}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config: r.RequiredCreate(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("id").MatchesOtherKey("name"),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("description").DoesNotExist(),
					check.That(data.ResourceName).Key("default_queue_type").HasValue("classic"),
					check.That(data.ResourceName).Key("max_connections").DoesNotExist(),
					check.That(data.ResourceName).Key("max_channels").DoesNotExist(),
					check.That(data.ResourceName).Key("tracing").HasValue("false"),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
			{
				Config: r.RequiredUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("id").MatchesOtherKey("name"),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("description").HasValue(r.Description),
					check.That(data.ResourceName).Key("default_queue_type").HasValue(r.DefaultQueueType),
					check.That(data.ResourceName).Key("max_connections").DoesNotExist(),
					check.That(data.ResourceName).Key("max_channels").DoesNotExist(),
					check.That(data.ResourceName).Key("tracing").HasValue("true"),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
		},
	})
}

func TestAccVhost_Optional(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_vhost", "test")
	r := acceptance.VhostResource{
		Name:             data.RandomString(),
		Description:      data.RandomString(),
		DefaultQueueType: "quorum",
		MaxConnections:   data.RandomIntegerString(),
		MaxQueues:        data.RandomIntegerString(),
		Tracing:          false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config: r.OptionalCreate(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("id").MatchesOtherKey("name"),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("description").HasValue(r.Description),
					check.That(data.ResourceName).Key("default_queue_type").HasValue(r.DefaultQueueType),
					check.That(data.ResourceName).Key("max_connections").HasValue(r.MaxConnections),
					check.That(data.ResourceName).Key("max_queues").HasValue(r.MaxQueues),
					check.That(data.ResourceName).Key("tracing").HasValue("false"),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
			{
				Config: r.OptionalUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("id").MatchesOtherKey("name"),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("description").HasValue(r.Description),
					check.That(data.ResourceName).Key("default_queue_type").HasValue(r.DefaultQueueType),
					check.That(data.ResourceName).Key("max_connections").HasValue(r.MaxConnections),
					check.That(data.ResourceName).Key("max_queues").HasValue(r.MaxQueues),
					check.That(data.ResourceName).Key("tracing").HasValue("true"),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
			{
				Config: r.OptionalUpdateLimits(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("id").MatchesOtherKey("name"),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("description").HasValue(r.Description),
					check.That(data.ResourceName).Key("default_queue_type").HasValue(r.DefaultQueueType),
					check.That(data.ResourceName).Key("max_connections").HasValue(r.MaxConnections),
					check.That(data.ResourceName).Key("max_queues").HasValue(r.MaxQueues),
					check.That(data.ResourceName).Key("tracing").HasValue("true"),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
			{
				Config: r.OptionalRemove(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("id").MatchesOtherKey("name"),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("description").RmqFeature(r.HasDescriptionUpdateFeature()).IsEmpty(),
					check.That(data.ResourceName).Key("default_queue_type").RmqFeature(r.HasDefaultQueueTypeUpdateFeature()).HasValue("classic"),
					check.That(data.ResourceName).Key("max_connections").IsEmpty(),
					check.That(data.ResourceName).Key("max_queues").IsEmpty(),
					check.That(data.ResourceName).Key("tracing").HasValue("false"),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
		},
	})
}

func TestAccVhost_ImportRequired(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_vhost", "test")
	r := acceptance.VhostResource{Name: data.RandomString()}

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
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: r.ImportStateVerifyIgnore(),
			},
		},
	})
}

func TestAccVhost_ImportOptional(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_vhost", "test")
	r := acceptance.VhostResource{
		Name:             data.RandomString(),
		Description:      data.RandomString(),
		DefaultQueueType: "stream",
		MaxConnections:   data.RandomIntegerString(),
		MaxQueues:        data.RandomIntegerString(),
		Tracing:          true,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config: r.OptionalCreate(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: r.ImportStateVerifyIgnore(),
			},
		},
	})
}

func TestAccVhost_AlredayExist(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_vhost", "test")
	r := acceptance.VhostResource{Name: "/"}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      r.RequiredCreate(data),
				ExpectError: regexp.MustCompile("vhost already exists"),
			},
		},
	})
}

func TestAccVhost_ErrorConvertingMaxConnections(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_vhost", "test")
	r := acceptance.VhostResource{Name: data.RandomString(), MaxConnections: data.RandomString()}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      r.ErrorConvertingCreate(data),
				ExpectError: regexp.MustCompile("error converting 'max_connections' to int"),
			},
			{
				Config: r.ErrorConvertingUpdate(data, data.RandomIntegerString(), r.MaxQueues),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("max_connections").HasValue(r.MaxConnections),
					check.That(data.ResourceName).Key("max_queues").IsEmpty(),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
			{
				Config:      r.ErrorConvertingUpdate(data, data.RandomString(), r.MaxQueues),
				ExpectError: regexp.MustCompile("error converting 'max_connections' to int"),
			},
		},
	})
}

func TestAccVhost_ErrorConvertingMaxQueues(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_vhost", "test")
	r := acceptance.VhostResource{Name: data.RandomString(), MaxQueues: data.RandomString()}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      r.ErrorConvertingCreate(data),
				ExpectError: regexp.MustCompile("error converting 'max_queues' to int"),
			},
			{
				Config: r.ErrorConvertingUpdate(data, r.MaxConnections, data.RandomIntegerString()),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("max_queues").HasValue(r.MaxQueues),
					check.That(data.ResourceName).Key("max_connections").IsEmpty(),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
			{
				Config:      r.ErrorConvertingUpdate(data, r.MaxConnections, data.RandomString()),
				ExpectError: regexp.MustCompile("error converting 'max_queues' to int"),
			},
		},
	})
}
