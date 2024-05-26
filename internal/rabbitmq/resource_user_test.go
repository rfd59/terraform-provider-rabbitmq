package rabbitmq_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/rfd59/terraform-provider-rabbitmq/internal/acceptance"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/acceptance/check"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUser_Required(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_user", "test")
	r := acceptance.UserResource{Name: data.RandomString(), Password: data.RandomString()}

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
					check.That(data.ResourceName).Key("password").HasValue(r.Password),
					check.That(data.ResourceName).Key("tags").DoesNotExist(),
					check.That(data.ResourceName).Key("max_connections").DoesNotExist(),
					check.That(data.ResourceName).Key("max_channels").DoesNotExist(),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
			{
				Config: r.RequiredUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("id").MatchesOtherKey("name"),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("password").HasValue(r.Password),
					check.That(data.ResourceName).Key("tags").DoesNotExist(),
					check.That(data.ResourceName).Key("max_connections").DoesNotExist(),
					check.That(data.ResourceName).Key("max_channels").DoesNotExist(),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
		},
	})
}

func TestAccUser_Optional(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_user", "test")
	r := acceptance.UserResource{
		Name:           data.RandomString(),
		Password:       data.RandomString(),
		Tags:           []string{"management"},
		MaxConnections: data.RandomIntegerString(),
		MaxChannels:    data.RandomIntegerString(),
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
					check.That(data.ResourceName).Key("password").HasValue(r.Password),
					check.That(data.ResourceName).Key("tags").Count(len(r.Tags)),
					check.That(data.ResourceName).Key("tags.0").HasValue(r.Tags[0]),
					check.That(data.ResourceName).Key("max_connections").HasValue(r.MaxConnections),
					check.That(data.ResourceName).Key("max_channels").HasValue(r.MaxChannels),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
			{
				Config: r.OptionalUpdateTags(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("id").MatchesOtherKey("name"),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("password").HasValue(r.Password),
					check.That(data.ResourceName).Key("tags").Count(len(r.Tags)),
					check.That(data.ResourceName).Key("tags.0").HasValue(r.Tags[0]),
					check.That(data.ResourceName).Key("max_connections").HasValue(r.MaxConnections),
					check.That(data.ResourceName).Key("max_channels").HasValue(r.MaxChannels),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
			{
				Config: r.OptionalUpdateLimits(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("id").MatchesOtherKey("name"),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("password").HasValue(r.Password),
					check.That(data.ResourceName).Key("tags").Count(len(r.Tags)),
					check.That(data.ResourceName).Key("tags.0").HasValue(r.Tags[0]),
					check.That(data.ResourceName).Key("max_connections").HasValue(r.MaxConnections),
					check.That(data.ResourceName).Key("max_channels").HasValue(r.MaxChannels),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
			{
				Config: r.OptionalRemove(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("id").MatchesOtherKey("name"),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("password").HasValue(r.Password),
					check.That(data.ResourceName).Key("tags").Count(len(r.Tags)),
					check.That(data.ResourceName).Key("max_connections").IsEmpty(),
					check.That(data.ResourceName).Key("max_channels").IsEmpty(),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
		},
	})
}

func TestAccUser_Login(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_user", "test")
	r := acceptance.UserResource{Name: data.RandomString(), Password: data.RandomString(), Tags: []string{"management"}}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config: r.LoginCreate(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("password").HasValue(r.Password),
					check.That(data.ResourceName).Key("tags").Count(len(r.Tags)),
					check.That(data.ResourceName).Key("tags.0").HasValue(r.Tags[0]),
					r.CheckLoginInRabbitMQ(),
				),
			},
			{
				Config: r.LoginUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("name").HasValue(r.Name),
					check.That(data.ResourceName).Key("password").HasValue(r.Password),
					check.That(data.ResourceName).Key("tags").Count(len(r.Tags)),
					check.That(data.ResourceName).Key("tags.0").HasValue(r.Tags[0]),
					r.CheckLoginInRabbitMQ(),
				),
			},
		},
	})
}

func TestAccUser_ImportRequired(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_user", "test")
	r := acceptance.UserResource{Name: data.RandomString(), Password: data.RandomString()}

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
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccUser_ImportOptional(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_user", "test")
	r := acceptance.UserResource{
		Name:           data.RandomString(),
		Password:       data.RandomString(),
		Tags:           []string{"management"},
		MaxConnections: data.RandomIntegerString(),
		MaxChannels:    data.RandomIntegerString(),
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
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccUser_AlredayExist(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_user", "test")
	r := acceptance.UserResource{Name: os.Getenv("RABBITMQ_USERNAME"), Password: os.Getenv("RABBITMQ_PASSWORD")}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      r.RequiredCreate(data),
				ExpectError: regexp.MustCompile("user already exists"),
			},
		},
	})
}

func TestAccUser_ErrorConvertingMaxConnections(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_user", "test")
	r := acceptance.UserResource{Name: data.RandomString(), Password: data.RandomString(), MaxConnections: data.RandomString()}

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
				Config: r.ErrorConvertingUpdate(data, data.RandomIntegerString(), r.MaxChannels),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("max_connections").HasValue(r.MaxConnections),
					check.That(data.ResourceName).Key("max_channels").IsEmpty(),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
			{
				Config:      r.ErrorConvertingUpdate(data, data.RandomString(), r.MaxChannels),
				ExpectError: regexp.MustCompile("error converting 'max_connections' to int"),
			},
		},
	})
}

func TestAccUser_ErrorConvertingMaxChannels(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_user", "test")
	r := acceptance.UserResource{Name: data.RandomString(), Password: data.RandomString(), MaxChannels: data.RandomString()}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      r.ErrorConvertingCreate(data),
				ExpectError: regexp.MustCompile("error converting 'max_channels' to int"),
			},
			{
				Config: r.ErrorConvertingUpdate(data, r.MaxConnections, data.RandomIntegerString()),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That(data.ResourceName).Key("max_channels").HasValue(r.MaxChannels),
					check.That(data.ResourceName).Key("max_connections").IsEmpty(),
					check.That(data.ResourceName).ExistsInRabbitMQ(r),
				),
			},
			{
				Config:      r.ErrorConvertingUpdate(data, r.MaxConnections, data.RandomString()),
				ExpectError: regexp.MustCompile("error converting 'max_channels' to int"),
			},
		},
	})
}
