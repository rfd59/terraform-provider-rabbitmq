package provider_test

import (
	"testing"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccQueue_importBasic(t *testing.T) {

	resourceName := "rabbitmq_queue.test"
	var queue rabbithole.QueueInfo

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: testAccQueueCheckDestroy(&queue),
		Steps: []resource.TestStep{
			{
				Config: testAccQueueConfig_basic,
				Check: testAccQueueCheck(
					resourceName, &queue,
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
