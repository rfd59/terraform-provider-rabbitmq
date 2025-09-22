package provider_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/rfd59/terraform-provider-rabbitmq/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBindings_DataSource(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	config := `
resource "rabbitmq_vhost" "test" {
	name = "accbindingsvhost"
}

resource "rabbitmq_exchange" "test" {
	name  = "accbindingsexchange"
	vhost = rabbitmq_vhost.test.name
	settings {
		type = "fanout"
	}
}

resource "rabbitmq_queue" "test" {
	name  = "accbindingsqueue"
	vhost = rabbitmq_vhost.test.name
	settings {
		durable = true
	}
}

resource "rabbitmq_binding" "test" {
	source           = rabbitmq_exchange.test.name
	vhost            = rabbitmq_vhost.test.name
	destination      = rabbitmq_queue.test.name
	destination_type = "queue"
	routing_key      = "#"
}

data "rabbitmq_bindings" "test" {
	vhost = rabbitmq_binding.test.vhost
}
`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAcc.PreCheck(t) },
		Providers: acceptance.TestAcc.Providers,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.rabbitmq_bindings.test", "bindings.0.source"),
					resource.TestCheckResourceAttrSet("data.rabbitmq_bindings.test", "bindings.0.destination"),
					resource.TestCheckResourceAttr("data.rabbitmq_bindings.test", "bindings.0.destination_type", "queue"),
					resource.TestCheckResourceAttr("data.rabbitmq_bindings.test", "bindings.0.routing_key", "#"),
					resource.TestCheckResourceAttr("data.rabbitmq_bindings.test", "vhost", "accbindingsvhost"),
				),
			},
		},
	})
}

func TestAccBindings_DataSourceNotExistVhost(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	data := acceptance.BuildTestData("rabbitmq_bindings", "test")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAcc.PreCheck(t) },
		Providers: acceptance.TestAcc.Providers,
		Steps: []resource.TestStep{
			{
				Config:      `data "` + data.ResourceType + `" "` + data.ResourceLabel + `" { vhost = "non-existent-vhost" }`,
				ExpectError: regexp.MustCompile(`Not Found`),
			},
		},
	})
}
