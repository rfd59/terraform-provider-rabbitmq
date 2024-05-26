package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/acceptance"
)

func TestAccVhost_importBasic(t *testing.T) {
	resourceName := "rabbitmq_vhost.test"
	var vhost string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: testAccVhostCheckDestroy(vhost),
		Steps: []resource.TestStep{
			{
				Config: testAccVhostConfig_basic,
				Check: testAccVhostCheck(
					resourceName, &vhost,
				),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
