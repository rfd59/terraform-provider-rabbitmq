package provider_test

import (
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/rfd59/terraform-provider-rabbitmq/internal/acceptance"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/acceptance/check"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExchange_DataSource(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	data := acceptance.BuildTestData("rabbitmq_exchange", "test")
	r := acceptance.ExchangeResource{Name: data.RandomString(), Vhost: "/", Settings: acceptance.ExchangeSettings{Type: "direct", Durable: true}}

	// Create an exchange to test the datasource
	r.SetDataSourceExchange(t)
	defer r.DelDataSourceExchange(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAcc.PreCheck(t) },
		Providers: acceptance.TestAcc.Providers,
		Steps: []resource.TestStep{
			{
				Config: r.DataSource(data),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+data.ResourceName).Exists(),
					check.That("data."+data.ResourceName).Key("name").HasValue(r.Name),
					check.That("data."+data.ResourceName).Key("vhost").HasValue(r.Vhost),
					check.That("data."+data.ResourceName).Key("id").HasValue(r.Name+"@"+r.Vhost),
					check.That("data."+data.ResourceName).Key("settings").Count(1),
					check.That("data."+data.ResourceName).Key("settings.0.type").HasValue(r.Settings.Type),
					check.That("data."+data.ResourceName).Key("settings.0.durable").HasValue(strconv.FormatBool(r.Settings.Durable)),
					check.That("data."+data.ResourceName).Key("settings.0.auto_delete").HasValue(strconv.FormatBool(r.Settings.AutoDelete)),
					check.That("data."+data.ResourceName).Key("settings.0.internal").HasValue(strconv.FormatBool(r.Settings.Internal)),
					check.That("data."+data.ResourceName).Key("settings.0.alternate_exchange").IsEmpty(),
					check.That("data."+data.ResourceName).Key("settings.0.arguments.%").DoesNotExist(),
				),
			},
		},
	})
}

func TestAccExchange_DataSourceNotExist(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_queue", "test")
	r := acceptance.QueueResource{Name: data.RandomString()}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAcc.PreCheck(t) },
		Providers: acceptance.TestAcc.Providers,
		Steps: []resource.TestStep{
			{
				Config:      r.DataSource(data),
				ExpectError: regexp.MustCompile("is not found"),
			},
		},
	})
}
