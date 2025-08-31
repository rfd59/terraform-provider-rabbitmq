package provider_test

import (
	"os"
	"regexp"
	"strconv"
	"testing"

	acceptance_test "github.com/rfd59/terraform-provider-rabbitmq/test/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExchangeFanout_DataSource(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	data := acceptance_test.BuildTestData("rabbitmq_exchange_fanout", "test")
	r := acceptance_test.ExchangeFanoutResource{
		ExchangeResource: acceptance_test.ExchangeResource{
			Name:              data.RandomString(),
			Vhost:             "/",
			Type:              "fanout",
			Durable:           true,
			AlternateExchange: data.RandomString(),
			Arguments: []map[string]interface{}{
				{"key": data.RandomString(), "value": data.RandomString(), "type": "string"},
			}}}

	// Create an exchange to test the datasource
	r.SetDataSourceExchange(t)
	//defer r.DelDataSourceExchange(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance_test.TestAcc.PreCheck(t) },
		Providers: acceptance_test.TestAcc.Providers,
		Steps: []resource.TestStep{
			{
				Config: r.DataSource(data),
				Check: resource.ComposeTestCheckFunc(
					acceptance_test.That("data."+data.ResourceName).Exists(),
					acceptance_test.That("data."+data.ResourceName).Key("id").HasValue(r.Name+"@"+r.Vhost),
					acceptance_test.That("data."+data.ResourceName).Key("name").HasValue(r.Name),
					acceptance_test.That("data."+data.ResourceName).Key("vhost").HasValue(r.Vhost),
					acceptance_test.That("data."+data.ResourceName).Key("type").HasValue(r.Type),
					acceptance_test.That("data."+data.ResourceName).Key("durable").HasValue(strconv.FormatBool(r.Durable)),
					acceptance_test.That("data."+data.ResourceName).Key("auto_delete").HasValue(strconv.FormatBool(r.AutoDelete)),
					acceptance_test.That("data."+data.ResourceName).Key("internal").HasValue(strconv.FormatBool(r.Internal)),
					acceptance_test.That("data."+data.ResourceName).Key("alternate_exchange").HasValue(r.AlternateExchange),
					acceptance_test.That("data."+data.ResourceName).Key("argument.#").Exists(),
					acceptance_test.That("data."+data.ResourceName).Key("argument.0.key").HasValue(r.Arguments[0]["key"].(string)),
					acceptance_test.That("data."+data.ResourceName).Key("argument.0.value").HasValue(r.Arguments[0]["value"].(string)),
					acceptance_test.That("data."+data.ResourceName).Key("argument.0.type").HasValue(r.Arguments[0]["type"].(string)),
				),
			},
		},
	})
}

func TestAccExchangeFanout_DataSourceNotExist(t *testing.T) {
	data := acceptance_test.BuildTestData("rabbitmq_exchange_fanout", "test")
	r := acceptance_test.ExchangeFanoutResource{
		ExchangeResource: acceptance_test.ExchangeResource{
			Name: data.RandomString(),
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance_test.TestAcc.PreCheck(t) },
		Providers: acceptance_test.TestAcc.Providers,
		Steps: []resource.TestStep{
			{
				Config:      r.DataSource(data),
				ExpectError: regexp.MustCompile("is not found"),
			},
		},
	})
}
