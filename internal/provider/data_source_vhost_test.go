package provider_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/rfd59/terraform-provider-rabbitmq/internal/acceptance"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/acceptance/check"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVhost_DataSource(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	data := acceptance.BuildTestData("rabbitmq_vhost", "test")
	r := acceptance.VhostResource{Name: data.RandomString()}

	// Create a vhost to test the datasource
	r.SetDataSourceVhost(t)

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
				),
			},
		},
	})

	// Remove the test vhost
	r.DelDataSourceVhost(t)
}

func TestAccVhost_DataSourceNotExist(t *testing.T) {
	data := acceptance.BuildTestData("rabbitmq_vhost", "test")
	r := acceptance.VhostResource{Name: data.RandomString()}

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
