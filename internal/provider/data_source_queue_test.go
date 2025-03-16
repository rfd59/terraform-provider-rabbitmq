package provider_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/rfd59/terraform-provider-rabbitmq/internal/acceptance"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/acceptance/check"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccQueue_DataSource(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	data := acceptance.BuildTestData("rabbitmq_queue", "test")
	r := acceptance.QueueResource{Name: data.RandomString(), Vhost: "/"}

	// Create a queue to test the datasource
	r.SetDataSourceQueue(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAcc.PreCheck(t) },
		Providers: acceptance.TestAcc.Providers,
		Steps: []resource.TestStep{
			{
				Config: r.DataSource(data),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+data.ResourceName).Exists(),
					check.That("data."+data.ResourceName).Key("id").MatchesRegex(regexp.MustCompile(r.Name+"@"+r.Vhost)),
					check.That("data."+data.ResourceName).Key("name").HasValue(r.Name),
					check.That("data."+data.ResourceName).Key("vhost").HasValue(r.Vhost),
					check.That("data."+data.ResourceName).Key("type").HasValue("classic"),
					check.That("data."+data.ResourceName).Key("status").IsNotEmpty(),
				),
			},
		},
	})

	// Remove the test queue
	r.DelDataSourceQueue(t)
}

func TestAccQueue_DataSourceNotExist(t *testing.T) {
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
