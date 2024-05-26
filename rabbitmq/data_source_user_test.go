package rabbitmq_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/rfd59/terraform-provider-rabbitmq/acceptance"
	"github.com/rfd59/terraform-provider-rabbitmq/acceptance/check"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUser_DataSource(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_user", "test")
	r := acceptance.UserResource{Name: os.Getenv("RABBITMQ_USERNAME")}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAcc.PreCheck(t) },
		Providers: acceptance.TestAcc.Providers,
		Steps: []resource.TestStep{
			{
				Config: r.DataSource(data),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+data.ResourceName).Exists(),
					check.That("data."+data.ResourceName).Key("id").MatchesOtherKey("name"),
					check.That("data."+data.ResourceName).Key("name").HasValue(r.Name),
					check.That("data."+data.ResourceName).Key("password").DoesNotExist(),
					check.That("data."+data.ResourceName).Key("tags").Count(1),
					check.That("data."+data.ResourceName).Key("max_connections").DoesNotExist(),
					check.That("data."+data.ResourceName).Key("max_channels").DoesNotExist(),
				),
			},
		},
	})
}

func TestAccUser_DataSourceLimits(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_user", "test")
	r := acceptance.UserResource{
		Name:           data.RandomString(),
		Password:       data.RandomString(),
		MaxConnections: data.RandomIntegerString(),
		MaxChannels:    data.RandomIntegerString(),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config: r.DataSourceLimits(data),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).Exists(),
					check.That("data."+data.ResourceName).Exists(),
					check.That("data."+data.ResourceName).Key("id").MatchesOtherKey("name"),
					check.That("data."+data.ResourceName).Key("name").HasValue(r.Name),
					check.That("data."+data.ResourceName).Key("password").DoesNotExist(),
					check.That("data."+data.ResourceName).Key("tags").DoesNotExist(),
					check.That("data."+data.ResourceName).Key("max_connections").HasValue(r.MaxConnections),
					check.That("data."+data.ResourceName).Key("max_channels").HasValue(r.MaxChannels),
				),
			},
		},
	})
}

func TestAccUser_DataSourceNotExist(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_user", "test")
	r := acceptance.UserResource{Name: data.RandomString()}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: r.CheckDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      r.DataSource(data),
				ExpectError: regexp.MustCompile("is not found"),
			},
		},
	})
}
